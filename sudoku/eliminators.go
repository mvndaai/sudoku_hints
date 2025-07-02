package sudoku

import (
	"cmp"
	"fmt"
	"slices"
)

var EliminatorFilledCell = CandidateEliminator{
	Name:        "Filled Cell",
	Description: "Eliminates candidates in the same row as a filled cell.",
	PartitionEliminator: func(cells []LocCell) (string, error) {
		found := []string{}
		for _, c := range cells {
			found = append(found, c.Cell.value)
		}
		for _, lc := range cells {
			removed := lc.Cell.RemoveCandiates(found)
			if len(removed) > 0 {
				return fmt.Sprintf("removed candidates (x:%d,y:%d) %v", lc.Loc.X, lc.Loc.Y, removed), nil
			}
		}
		return "", nil
	},
	Simple: true,
}

var EliminatorMatchingCandidates = CandidateEliminator{
	Name:        "Matching Candidates",
	Description: "Eliminates candidates if any cells have the same candidates.",
	PartitionEliminator: func(cells []LocCell) (string, error) {
		candidatesCount := map[string]int{}
		candidatesValues := map[string][]string{}

		for _, c := range cells {
			if len(c.Cell.Candidates) == 0 {
				continue // Skip empty candidates
			}
			cs := fmt.Sprint(c.Cell.Candidates)
			candidatesCount[cs]++
			candidatesValues[cs] = c.Cell.Candidates
		}

		var viableGroups [][]string
		for cs, count := range candidatesCount {
			if count < 2 {
				continue // We need at least two cells with the same candidates
			}
			if len(candidatesValues[cs]) != count {
				continue // Ensure the candidates the number of candiates matches the count
			}
			viableGroups = append(viableGroups, candidatesValues[cs])
		}

		for _, vg := range viableGroups {
			for _, lc := range cells {
				if slices.Equal(lc.Cell.Candidates, vg) {
					continue
				}
				removed := lc.Cell.RemoveCandiates(vg)
				if len(removed) > 0 {
					return fmt.Sprintf("removed candidates (x:%d,y:%d) %v", lc.Loc.X, lc.Loc.Y, removed), nil
				}
			}
		}
		return "", nil
	},
}

var EliminatorGroupAndRowColumn = CandidateEliminator{
	Name:        "Group and Row/Column",
	Description: "If a group only has values in a row or column, then remove those candidates from the other cells in that row or column.",
	GameEliminator: func(g *Game) (string, error) {
		groupPlaceValues := map[int]map[string]map[int][]string{}

		var rowKey, colKey = "row", "col"

		for y, row := range g.Board {
			for x, gc := range row {
				for _, c := range gc.Cell.Candidates {
					if groupPlaceValues[gc.group] == nil {
						groupPlaceValues[gc.group] = map[string]map[int][]string{}
					}
					if groupPlaceValues[gc.group][colKey] == nil {
						groupPlaceValues[gc.group][colKey] = map[int][]string{}
					}
					if _, exists := groupPlaceValues[gc.group][colKey][x]; !exists {
						groupPlaceValues[gc.group][colKey][x] = []string{}
					}
					groupPlaceValues[gc.group][colKey][x] = append(groupPlaceValues[gc.group][colKey][x], c)

					if groupPlaceValues[gc.group][rowKey] == nil {
						groupPlaceValues[gc.group][rowKey] = map[int][]string{}
					}
					if _, exists := groupPlaceValues[gc.group][rowKey][y]; !exists {
						groupPlaceValues[gc.group][rowKey][y] = []string{}
					}
					groupPlaceValues[gc.group][rowKey][y] = append(groupPlaceValues[gc.group][rowKey][y], c)

					slices.Sort(groupPlaceValues[gc.group][colKey][x])
					groupPlaceValues[gc.group][colKey][x] = slices.Compact(groupPlaceValues[gc.group][colKey][x])
					slices.Sort(groupPlaceValues[gc.group][rowKey][y])
					groupPlaceValues[gc.group][rowKey][y] = slices.Compact(groupPlaceValues[gc.group][rowKey][y])
				}
			}
		}

		// Check if any group only has a candidate in a single row or column
		groupPlaceLocs := map[int]map[string]map[string][]int{}
		for group, places := range groupPlaceValues {
			for placeType, locs := range places {
				for loc, candidates := range locs {
					for _, c := range candidates {
						if groupPlaceLocs[group] == nil {
							groupPlaceLocs[group] = map[string]map[string][]int{}
						}
						if groupPlaceLocs[group][placeType] == nil {
							groupPlaceLocs[group][placeType] = map[string][]int{}
						}
						if groupPlaceLocs[group][placeType][c] == nil {
							groupPlaceLocs[group][placeType][c] = []int{}
						}
						groupPlaceLocs[group][placeType][c] = append(groupPlaceLocs[group][placeType][c], loc)
					}
				}
			}
		}

		type holder struct {
			group int
			loc   int
		}
		rowsToRemove := map[holder][]string{}
		colsToRemove := map[holder][]string{}

		//holders := map[holder][]string{}
		for group, places := range groupPlaceLocs {
			for placeType, locs := range places {
				for candiate, rows := range locs {
					if len(rows) != 1 {
						continue
					}
					holder := holder{group: group, loc: rows[0]}
					if placeType == rowKey {
						if _, exists := rowsToRemove[holder]; !exists {
							rowsToRemove[holder] = []string{}
						}
						rowsToRemove[holder] = append(rowsToRemove[holder], candiate)
						continue
					}
					if placeType == colKey {
						if _, exists := colsToRemove[holder]; !exists {
							colsToRemove[holder] = []string{}
						}
						colsToRemove[holder] = append(colsToRemove[holder], candiate)
						continue
					}
				}
			}
		}

		if len(rowsToRemove) == 0 && len(colsToRemove) == 0 {
			return "", nil // No candidates to remove
		}

		type candiateHolder struct {
			group     int
			loc       int
			candiates []string
		}
		rowsToRemoveSlice := make([]candiateHolder, 0, len(rowsToRemove))
		colsToRemoveSlice := make([]candiateHolder, 0, len(colsToRemove))
		for h, c := range rowsToRemove {

			rowsToRemoveSlice = append(rowsToRemoveSlice, candiateHolder{group: h.group, loc: h.loc, candiates: c})
		}
		for h, c := range colsToRemove {

			colsToRemoveSlice = append(colsToRemoveSlice, candiateHolder{group: h.group, loc: h.loc, candiates: c})

		}

		sortFunc := func(a, b candiateHolder) int {
			if a.group != b.group {
				return cmp.Compare(a.group, b.group)
			}
			if a.loc != b.loc {
				return cmp.Compare(a.loc, b.loc)
			}
			return slices.Compare(a.candiates, b.candiates)
		}
		slices.SortFunc(rowsToRemoveSlice, sortFunc)
		slices.SortFunc(colsToRemoveSlice, sortFunc)

		for y, row := range g.Board {
			for x, gc := range row {
				i := slices.IndexFunc(rowsToRemoveSlice, func(h candiateHolder) bool {
					return y == h.loc && gc.group != h.group
				})
				if i >= 0 {
					removed := gc.Cell.RemoveCandiates(rowsToRemoveSlice[i].candiates)
					if len(removed) > 0 {
						return fmt.Sprintf("removed candidates (x:%d,y:%d) %v", x, y, removed), nil
					}
				}
				i = slices.IndexFunc(colsToRemoveSlice, func(h candiateHolder) bool {
					return x == h.loc && gc.group != h.group
				})
				if i >= 0 {
					removed := gc.Cell.RemoveCandiates(colsToRemoveSlice[i].candiates)
					if len(removed) > 0 {
						return fmt.Sprintf("removed candidates (x:%d,y:%d) %v", x, y, removed), nil
					}
				}
			}
		}
		return "", nil
	},
}

/* Rules to add
** If cells in the a partition group out the same candidates remove those from the others. Example [1 4], [4 6], [1 6] or [1 4], [4 6], [1 4 7].
** If a group only has values in a row or column, then remove those candidates from the other cells in that row or column.
 */

package sudoku

import (
	"cmp"
	"fmt"
	"slices"
)

func PartitionHinterToEliminator(bpe PartitionHinter) PartitionEliminator {
	return func(cells []LocCell) (string, error) {
		ok, h, err := bpe(cells)
		if err != nil {
			return "", err
		}
		if !ok {
			return "", nil
		}
		_ = h.cell.RemoveCandiates(h.CandidatesToRemove)
		return fmt.Sprintf("removed candidates (x:%d,y:%d) %v", h.Loc.X, h.Loc.Y, h.CandidatesToRemove), nil
	}
}

var EliminatorFilledCell = func() CandidateEliminator {
	name := "Filled Cell"
	r := CandidateEliminator{
		Name:        name,
		Description: "Eliminates candidates in the same row as a filled cell.",
		PartitionHinter: func(cells []LocCell) (bool, Hint, error) {
			found := []string{}
			for _, c := range cells {
				found = append(found, c.Cell.Value)
			}
			for _, lc := range cells {
				diffs := lc.Cell.CandidateDiffs(found)
				if len(diffs) > 0 {
					return true, Hint{
						Loc:                lc.Loc,
						Eliminator:         name,
						CandidatesToRemove: diffs,
						cell:               lc.Cell,
					}, nil
				}
			}
			return false, Hint{}, nil
		},
		Simple: true,
	}

	r.PartitionEliminator = PartitionHinterToEliminator(r.PartitionHinter)
	return r
}()

type Locs []Loc

func (l Locs) Key() string {
	// Sort the locations to ensure consistent key generation
	slices.SortFunc(l, func(a, b Loc) int {
		if a.Y != b.Y {
			return cmp.Compare(a.Y, b.Y)
		}
		return cmp.Compare(a.X, b.X)
	})
	return fmt.Sprint(l)
}

// TODO expand as hidden pairs/triples
var EliminatorUniqueCandidate = CandidateEliminator{
	Name:        "Unique Candidate",
	Description: "Eliminates all other candidates if a cell has a unique candidate in its partition.",
	PartitionEliminator: func(cells []LocCell) (string, error) {
		candidates := map[string]Locs{}
		for _, c := range cells {
			for _, candidate := range c.Cell.Candidates {
				if _, exists := candidates[candidate]; !exists {
					candidates[candidate] = []Loc{}
				}
				candidates[candidate] = append(candidates[candidate], c.Loc)
			}
		}

		type uniqueCandidate struct {
			candiates []string
			locs      Locs
		}
		uniqueCandidatesBuilder := map[string]uniqueCandidate{}
		for candidate, locs := range candidates {
			key := locs.Key()
			uc, exists := uniqueCandidatesBuilder[key]
			if !exists {
				uc = uniqueCandidate{candiates: []string{}, locs: locs}
			}
			uc.candiates = append(uc.candiates, candidate)
			uniqueCandidatesBuilder[key] = uc
		}

		// Remove ones where len of candidates does not match the number of locations
		for key, uc := range uniqueCandidatesBuilder {
			if len(uc.candiates) != len(uc.locs) {
				delete(uniqueCandidatesBuilder, key)
				continue
			}
		}

		if len(uniqueCandidatesBuilder) == 0 {
			return "", nil // No unique candidates found
		}

		uniqueCandidates := map[Loc][]string{}
		for _, uc := range uniqueCandidatesBuilder {
			for _, loc := range uc.locs {
				uniqueCandidates[loc] = uc.candiates
			}
		}

		for _, lc := range cells {
			ucs, ok := uniqueCandidates[lc.Loc]
			if !ok {
				continue // No unique candidate for this cell
			}
			toRemove := slices.DeleteFunc(slices.Clone(lc.Cell.Candidates), func(c string) bool {
				return slices.Contains(ucs, c)
			})
			removed := lc.Cell.RemoveCandiates(toRemove)
			if len(removed) > 0 {
				return fmt.Sprintf("removed candidates (x:%d,y:%d) %v", lc.Loc.X, lc.Loc.Y, removed), nil
			}
		}
		return "", nil
	},
	Simple: true,
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

var EliminatorCandidateChains = CandidateEliminator{
	Name:        "Candidate Chains",
	Description: "If N cells form a chain where they share exactly N candidates total, remove those candidates from other cells in the partition.",
	PartitionEliminator: func(cells []LocCell) (string, error) {
		// Get cells with candidates
		candidateCells := []LocCell{}
		for _, lc := range cells {
			if len(lc.Cell.Candidates) >= 2 {
				candidateCells = append(candidateCells, lc)
			}
		}

		// Try all possible chain sizes from 2 to the number of cells
		for chainSize := 2; chainSize <= len(candidateCells); chainSize++ {
			// Generate all combinations of cells of the given chain size
			combinations := getCombinations(candidateCells, chainSize)

			for _, combo := range combinations {
				// Collect all unique candidates from the cells in this combination
				allCandidates := map[string]bool{}
				for _, cell := range combo {
					for _, c := range cell.Cell.Candidates {
						allCandidates[c] = true
					}
				}

				// Check if this forms a valid chain (N cells with N total candidates)
				if len(allCandidates) == chainSize {
					candidatesList := make([]string, 0, len(allCandidates))
					for c := range allCandidates {
						candidatesList = append(candidatesList, c)
					}

					// Create a map of locations in the chain for quick lookup
					chainLocs := map[Loc]bool{}
					for _, cell := range combo {
						chainLocs[cell.Loc] = true
					}

					// Remove these candidates from all other cells
					for _, lc := range cells {
						if chainLocs[lc.Loc] {
							continue
						}
						removed := lc.Cell.RemoveCandiates(candidatesList)
						if len(removed) > 0 {
							return fmt.Sprintf("removed candidates (x:%d,y:%d) %v from chain of size %d", lc.Loc.X, lc.Loc.Y, removed, chainSize), nil
						}
					}
				}
			}
		}
		return "", nil
	},
}

// Helper function to generate all combinations of a given size
func getCombinations(cells []LocCell, size int) [][]LocCell {
	if size == 0 {
		return [][]LocCell{{}}
	}
	if len(cells) < size {
		return [][]LocCell{}
	}

	result := [][]LocCell{}

	// Include first element
	for _, combo := range getCombinations(cells[1:], size-1) {
		newCombo := make([]LocCell, len(combo)+1)
		newCombo[0] = cells[0]
		copy(newCombo[1:], combo)
		result = append(result, newCombo)
	}

	// Exclude first element
	result = append(result, getCombinations(cells[1:], size)...)

	return result
}

// EliminatorFistemafelRing enforces that certain groups must contain the same set of values
// This is used for variant Sudoku where groups are "disjoint" - they must have matching values
// https://www.tiktok.com/@brainfueltips/video/7565584522092268813
var EliminatorFistemafelRing = CandidateEliminator{
	Name:        "Fistemafel Ring",
	Description: "The 16 digits that ring the center much match the corners.",
	GameEliminator: func(g *Game) (string, error) {
		// Define the specific cells for each matching group by their coordinates
		matchingGroups := []struct {
			name string
			locs []Loc
		}{
			{
				name: "corners",
				locs: []Loc{
					{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 7, Y: 0}, {X: 8, Y: 0},
					{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 7, Y: 1}, {X: 8, Y: 1},
					{X: 0, Y: 7}, {X: 1, Y: 7}, {X: 7, Y: 7}, {X: 8, Y: 7},
					{X: 0, Y: 8}, {X: 1, Y: 8}, {X: 7, Y: 8}, {X: 8, Y: 8},
				},
			},
			{
				name: "ring",
				locs: []Loc{
					{X: 2, Y: 2}, {X: 3, Y: 2}, {X: 4, Y: 2}, {X: 5, Y: 2}, {X: 6, Y: 2},
					{X: 2, Y: 3}, {X: 6, Y: 3},
					{X: 2, Y: 4}, {X: 6, Y: 4},
					{X: 2, Y: 5}, {X: 6, Y: 5},
					{X: 2, Y: 6}, {X: 3, Y: 6}, {X: 4, Y: 6}, {X: 5, Y: 6}, {X: 6, Y: 6},
				},
			},
		}

		// Collect cells for each group
		groupCells := make([][]LocCell, len(matchingGroups))
		for i, group := range matchingGroups {
			groupCells[i] = make([]LocCell, 0, len(group.locs))
			for _, loc := range group.locs {
				if loc.Y < len(g.Board) && loc.X < len(g.Board[loc.Y]) {
					groupCells[i] = append(groupCells[i], LocCell{
						Loc:  loc,
						Cell: g.Board[loc.Y][loc.X].Cell,
					})
				}
			}
		}

		// Check which groups are complete
		groupValues := make([]map[string]bool, len(groupCells))
		groupComplete := make([]bool, len(groupCells))

		for i, cells := range groupCells {
			groupValues[i] = make(map[string]bool)
			filledCount := 0

			for _, cell := range cells {
				if cell.Cell.Value != "" {
					groupValues[i][cell.Cell.Value] = true
					filledCount++
				}
			}

			groupComplete[i] = (filledCount == len(cells))
		}

		// If exactly one group is complete, remove candidates from other groups that aren't in the complete group's values
		completeGroupIdx := -1
		completeCount := 0

		for i, complete := range groupComplete {
			if complete {
				completeGroupIdx = i
				completeCount++
			}
		}

		if completeCount != 1 {
			return "", nil
		}

		// One group is complete - enforce its values on other groups
		completeValues := groupValues[completeGroupIdx]
		completeGroupName := matchingGroups[completeGroupIdx].name

		for i, cells := range groupCells {
			if i == completeGroupIdx {
				continue // Skip the complete group
			}

			groupName := matchingGroups[i].name

			// Remove candidates that are not in the complete group's values
			for _, cell := range cells {
				if cell.Cell.Value != "" {
					continue
				}

				toRemove := []string{}
				for _, candidate := range cell.Cell.Candidates {
					if !completeValues[candidate] {
						toRemove = append(toRemove, candidate)
					}
				}

				if len(toRemove) == 0 {
					continue
				}

				removed := cell.Cell.RemoveCandiates(toRemove)
				if len(removed) > 0 {
					return fmt.Sprintf("removed %v from %s cell at (%d,%d) because %s is complete without these values", removed, groupName, cell.Loc.X, cell.Loc.Y, completeGroupName), nil
				}
			}
		}

		return "", nil
	},
}

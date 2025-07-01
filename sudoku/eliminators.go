package sudoku

import (
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

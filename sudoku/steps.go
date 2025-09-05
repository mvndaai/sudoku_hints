package sudoku

import (
	"fmt"
	"log"
	"slices"
)

func (g *Game) RemoveAllSimple(clearRecentCandidates bool) error {
	log.Println("in RemoveAllSimple")
	eliminator := EliminatorFilledCell
	rows, cols, groups := g.GetSectionedCells()
	partitions := []struct {
		name  string
		cells [][]LocCell
	}{
		{"rows", rows},
		{"cols", cols},
		{"groups", groups},
	}

	for _, ps := range partitions {
		for i := range ps.cells {
			for {
				ok, hint, err := eliminator.PartitionHinter(ps.cells[i])
				if err != nil {
					return fmt.Errorf("(%s) %s %d: %w", hint.Eliminator, ps.name, i, err)
				}
				if !ok {
					//log.Println(cells, "no change")
					break
				}

				bfCandidates := slices.Clone(hint.cell.Candidates)
				bfRecentCandidates := slices.Clone(hint.cell.RecentCandidates)
				removed := hint.cell.RemoveCandiates(hint.CandidatesToRemove)
				if len(removed) != 0 {
					log.Println("Removed candidates:", hint.CandidatesToRemove, bfCandidates, bfRecentCandidates)
					log.Println("After:", hint.cell.Candidates, hint.cell.RecentCandidates)
				}
			}
		}
	}

	if clearRecentCandidates {
		g.RemoveAllRecentCandidates()
	}
	return nil
}

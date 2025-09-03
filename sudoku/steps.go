package sudoku

import (
	"fmt"
)

func (g *Game) RemoveAllSimple(clearRecentCandidates bool) error {
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
				ok, hint, err := eliminator.BetterPartitionEliminator(ps.cells[i])
				if err != nil {
					return fmt.Errorf("(%s) %s %d: %w", hint.Eliminator, ps.name, i, err)
				}
				if !ok {
					//log.Println(cells, "no change")
					break
				}

				//log.Println("Removed candidates:", hint.CandidatesToRemove)
				_ = hint.cell.RemoveCandiates(hint.CandidatesToRemove)
				//log.Println("After:", hint.cell.RecentCandidates)
			}
		}
	}

	//if clearRecentCandidates {
	//	g.RemoveAllRecentCandidates()
	//}
	return nil
}

func NextHint(g *Game, eliminators []string) (bool, Hint, error) {

	// Get the next hint from the game state
	return false, Hint{}, nil
}

package sudoku

import (
	"fmt"
	"log"
)

func (g *Game) RemoveAllSimple() error {
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
		for i, cells := range ps.cells {
			for {
				ok, hint, err := eliminator.BetterPartitionEliminator(cells)
				if err != nil {
					return fmt.Errorf("(%s) %s %d: %w", hint.Eliminator, ps.name, i, err)
				}
				if !ok {
					log.Println(cells, "no change")
					break
				}
				_ = hint.cell.RemoveCandiates(hint.CandidatesToRemove)
			}
		}
	}
	return nil
}

func NextHint(g *Game, eliminators []string) (bool, Hint, error) {

	// Get the next hint from the game state
	return false, Hint{}, nil
}

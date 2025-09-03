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
				ok, hint, err := eliminator.PartitionHinter(ps.cells[i])
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

func NextHint(g *Game, eliminators []string) (Hint, error) {
	rows, cols, groups := g.GetSectionedCells()

	for _, eliminator := range Eliminators {
		if eliminator.PartitionEliminator != nil {
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
					change, err := eliminator.PartitionEliminator(cells)
					if err != nil {
						//return "", fmt.Errorf("(%s) %s %d: %w", eliminator.Name, ps.name, i, err)
						return Hint{}, fmt.Errorf("(%s) %s %d: %w", eliminator.Name, ps.name, i, err)
					}
					if change != "" {
						//return fmt.Sprintf("(%s) %s %d: %s", eliminator.Name, ps.name, i, change), nil
					}
				}
			}
		}
		if eliminator.GameEliminator != nil {
			change, err := eliminator.GameEliminator(g)
			if err != nil {
				return Hint{}, fmt.Errorf("(%s): %w", eliminator.Name, err)
			}
			if change != "" {
				//return fmt.Sprintf("(%s): %s", eliminator.Name, change), nil
				return Hint{}, nil
			}
		}
	}
	return Hint{}, fmt.Errorf("no candidates eliminated by any rules")
}

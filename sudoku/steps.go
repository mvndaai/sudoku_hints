package sudoku

import (
	"fmt"
	"log"
	"slices"
)

func (g *Game) RemoveAllSimple(clearRecentCandidates bool) error {
	//log.Println("in RemoveAllSimple")
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

func (g *Game) RemoveOneCandidate(clearRecentCandidates bool) (bool, string, error) {
	if clearRecentCandidates {
		g.RemoveAllRecentCandidates()
	}

	// First check if there are any candidates at all
	candidateCount := 0
	for _, row := range g.Board {
		for _, gc := range row {
			candidateCount += len(gc.Cell.Candidates)
		}
	}
	//log.Printf("RemoveOneCandidate: Total candidates on board: %d\n", candidateCount)

	//log.Println("RemoveOneCandidate: Starting search for candidates to remove")

	rows, cols, groups := g.GetSectionedCells()

	// Iterate through all eliminators
	for _, eliminator := range Eliminators {
		// Try PartitionEliminator if it exists
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
				for i := range ps.cells {
					// Track candidates before elimination
					candidatesBeforeMap := make(map[Loc][]string)
					for _, lc := range ps.cells[i] {
						if len(lc.Cell.Candidates) > 0 {
							candidatesBeforeMap[lc.Loc] = slices.Clone(lc.Cell.Candidates)
						}
					}

					change, err := eliminator.PartitionEliminator(ps.cells[i])
					if err != nil {
						return false, "", fmt.Errorf("(%s) %s %d: %w", eliminator.Name, ps.name, i, err)
					}

					if change != "" {
						// Check if any candidates were actually removed
						for _, lc := range ps.cells[i] {
							before, existed := candidatesBeforeMap[lc.Loc]
							if existed && len(before) != len(lc.Cell.Candidates) {
								log.Printf("RemoveOneCandidate: Successfully removed candidates using %s: %s\n", eliminator.Name, change)
								changeWithName := fmt.Sprintf("[%s] %s %d: %s", eliminator.Name, ps.name, i, change)
								return true, changeWithName, nil
							}
						}
					}
				}
			}
		}

		// Try GameEliminator if it exists
		if eliminator.GameEliminator != nil {
			// Track total candidates before
			totalBefore := 0
			for _, row := range g.Board {
				for _, gc := range row {
					totalBefore += len(gc.Cell.Candidates)
				}
			}

			change, err := eliminator.GameEliminator(g)
			if err != nil {
				return false, "", fmt.Errorf("(%s): %w", eliminator.Name, err)
			}

			if change != "" {
				// Check if any candidates were removed
				totalAfter := 0
				for _, row := range g.Board {
					for _, gc := range row {
						totalAfter += len(gc.Cell.Candidates)
					}
				}

				if totalAfter < totalBefore {
					log.Printf("RemoveOneCandidate: Successfully removed candidates using %s: %s\n", eliminator.Name, change)
					changeWithName := fmt.Sprintf("[%s] %s", eliminator.Name, change)
					return true, changeWithName, nil
				}
			}
		}
	}

	//log.Println("RemoveOneCandidate: No candidates found to remove across all eliminators")
	return false, "", nil
}

package sudoku

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

type (
	PartitionEliminator func([]LocCell) (change string, _ error)
	GameEliminator      func(*Game) (change string, _ error)
)

type CandidateEliminator struct {
	Name                string
	Description         string
	PartitionEliminator PartitionEliminator
	GameEliminator      GameEliminator
	Simple              bool // Allow hiding simple eliminators from the UI
}

// Rules is a collection of rules that can be applied to a Sudoku game.
var Eliminators = []CandidateEliminator{
	EliminatorFilledCell,
	EliminatorUniqueCandidate,
	EliminatorMatchingCandidates,
	EliminatorGroupAndRowColumn,
}

func (g *Game) GetSectionedCells() (rows [][]LocCell, cols [][]LocCell, groups [][]LocCell) {
	rows = make([][]LocCell, len(g.Symbols))
	cols = make([][]LocCell, len(g.Symbols))
	groupMap := make(map[int][]LocCell)
	for y := range g.Board {
		rows[y] = make([]LocCell, 0, len(g.Board[y]))
		for x := range g.Board[y] {
			lc := LocCell{
				Loc:  Loc{X: x, Y: y},
				Cell: g.Board[y][x].Cell,
			}
			rows[y] = append(rows[y], lc)
			if cols[x] == nil {
				cols[x] = make([]LocCell, 0, len(g.Board))
			}
			cols[x] = append(cols[x], lc)

			if _, exists := groupMap[g.Board[y][x].group]; !exists {
				groupMap[g.Board[y][x].group] = []LocCell{}
			}
			groupMap[g.Board[y][x].group] = append(groupMap[g.Board[y][x].group], lc)
		}
	}

	groups = make([][]LocCell, len(g.Symbols))
	for i, cells := range groupMap {
		groups[i] = cells
	}
	return rows, cols, groups
}

func (g *Game) EliminateCandidates(onlySimples bool) (change string, _ error) {
	rows, cols, groups := g.GetSectionedCells()

	if g.RandomEliminators {
		// Shuffle the eliminators to randomize the order of elimination

		slices.SortFunc(Eliminators, func(a, b CandidateEliminator) int {
			if a.Simple != b.Simple {
				if a.Simple {
					return -1 // Simple eliminators come first
				}
				return 1 // Non-simple eliminators come after simple ones
			}
			return rand.IntN(2)*2 - 1 // Randomly order them if they are both simple or both not simple
		})
	}
	for _, eliminator := range Eliminators {
		if onlySimples && !eliminator.Simple {
			continue // Skip non-simple eliminators if onlySimples is true
		}
		if eliminator.PartitionEliminator != nil {
			partitions := []struct {
				name  string
				cells [][]LocCell
			}{
				{"rows", rows},
				{"cols", cols},
				{"groups", groups},
			}

			if g.RandomEliminators {
				// Shuffle partitions to randomize the order of elimination
				rand.Shuffle(len(partitions), func(i, j int) { partitions[i], partitions[j] = partitions[j], partitions[i] })
			}

			for _, ps := range partitions {
				for i, cells := range ps.cells {
					change, err := eliminator.PartitionEliminator(cells)
					if err != nil {
						return "", fmt.Errorf("(%s) %s %d: %w", eliminator.Name, ps.name, i, err)
					}
					if change != "" {
						if g.HideSimple && eliminator.Simple {
							return "", nil // Skip basic eliminators if HideBasic is true
						}
						return fmt.Sprintf("(%s) %s %d: %s", eliminator.Name, ps.name, i, change), nil
					}
				}
			}
		}
		if eliminator.GameEliminator != nil {
			change, err := eliminator.GameEliminator(g)
			if err != nil {
				return "", fmt.Errorf("(%s): %w", eliminator.Name, err)
			}
			if change != "" {
				return fmt.Sprintf("(%s): %s", eliminator.Name, change), nil
			}
		}
	}
	return "", fmt.Errorf("no candidates eliminated by any rules")
}

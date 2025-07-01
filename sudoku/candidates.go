package sudoku

import (
	"fmt"
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
}

// Rules is a collection of rules that can be applied to a Sudoku game.
var Eliminators = []CandidateEliminator{
	EliminatorFilledCell,
	EliminatorMatchingCandidates,
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

func (g *Game) EliminateCandidates() (change string, _ error) {
	rows, cols, groups := g.GetSectionedCells()

	for _, eliminator := range Eliminators {
		//fmt.Println("Running eliminator:", eliminator.Name)

		if eliminator.PartitionEliminator != nil {
			for i, row := range rows {
				change, err := eliminator.PartitionEliminator(row)
				if err != nil {
					return "", fmt.Errorf("(%s )row %d: %w", eliminator.Name, i, err)
				}
				if change != "" {
					return fmt.Sprintf("(%s) row %d: %s", eliminator.Name, i, change), nil
				}
			}
			for i, col := range cols {
				change, err := eliminator.PartitionEliminator(col)
				if err != nil {
					return "", fmt.Errorf("(%s) col %d: %w", eliminator.Name, i, err)
				}
				if change != "" {
					return fmt.Sprintf("(%s) col %d: %s", eliminator.Name, i, change), nil
				}
			}
			for i, group := range groups {
				change, err := eliminator.PartitionEliminator(group)
				if err != nil {
					return "", fmt.Errorf("(%s) group %d: %w", eliminator.Name, i, err)
				}
				if change != "" {
					return fmt.Sprintf("(%s) group %d: %s", eliminator.Name, i, change), nil
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

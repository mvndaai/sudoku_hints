package sudoku

type CandidateEliminator struct {
	Name        string
	Description string
	Executor    func([]*Cell) error // This takes a slice of Cells because it extracts
}

// Rules is a collection of rules that can be applied to a Sudoku game.
var Eliminators = []CandidateEliminator{
	{
		Name:        "Filled Cell",
		Description: "Eliminates candidates in the same row as a filled cell.",
		Executor: func(cells []*Cell) error {
			found := []string{}
			for _, c := range cells {
				found = append(found, c.value)
			}
			for i := range cells {
				cells[i].RemoveCandiates(found)
			}
			return nil
		},
	},
}

// TODO I might need to have loc to show location for hints
func (g *Game) GetSectionedCells() (rows [][]*Cell, cols [][]*Cell, groups [][]*Cell) {
	rows = make([][]*Cell, 0, len(g.Symbols))
	cols = make([][]*Cell, 0, len(g.Symbols))
	groupMap := make(map[int][]*Cell)
	for y := range g.Board {
		for x := range g.Board[y] {
			cell := &g.Board[y][x].Cell
			rows[y] = append(rows[y], cell)
			cols[x] = append(cols[x], cell)

			if _, exists := groupMap[g.Board[y][x].group]; !exists {
				groupMap[g.Board[y][x].group] = []*Cell{}
			}
			groupMap[g.Board[y][x].group] = append(groupMap[g.Board[y][x].group], cell)
		}
	}

	groups = make([][]*Cell, len(g.Symbols))
	for i, cells := range groupMap {
		groups[i] = cells
	}
	return rows, cols, groups
}

func (g *Game) EliminateCandidates() error {
	//rows, cols, groups := g.GetSectionedCells()

	//for _, eliminator := range Eliminators {
	//	if err := eliminator.Executor(rows); err != nil {
	//		return err
	//	}
	//	if err := eliminator.Executor(cols); err != nil {
	//		return err
	//	}
	//	if err := eliminator.Executor(groups); err != nil {
	//		return err
	//	}
	//}

	return nil

}

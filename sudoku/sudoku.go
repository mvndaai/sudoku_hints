package sudoku

import (
	"fmt"
	"slices"
	"strconv"
)

type (
	GroupedCell struct {
		group int
		Cell  *Cell `json:"cell"`
	}

	LocCell struct {
		Loc  Loc
		Cell *Cell
	}

	Cell struct {
		Value        string   `json:"value"`
		Candidates   []string `json:"candidates"`
		IsPreFilled  bool     `json:"isPreFilled"` // If true, this cell was part of the original puzzle and should not be changed
		IsLastFilled bool     `json:"isLastFilled"`
	}

	Game struct {
		Symbols           []string
		Board             [][]GroupedCell
		HideSimple        bool
		RandomEliminators bool // If true, the eliminators will be run in a random order
		RunSimpleFirst    bool // If true, the simple eliminators will be run quietly first
		RunOnce           bool // If true, breaks after finding one value
		AutoSolve         bool
	}
)

func (g *Game) Fill(cells [][]string, group map[Loc]int, symbols []string) error {
	symbolMap := map[string]struct{}{}

	// Fill in the cell values and track unique symbols
	g.Board = make([][]GroupedCell, len(cells))
	for y, row := range cells {
		g.Board[y] = make([]GroupedCell, len(row))
		for x, v := range row {
			var hasStartingValue bool
			if v != "" {
				symbolMap[v] = struct{}{}
				hasStartingValue = true
			}
			g.Board[y][x] = GroupedCell{
				Cell:  &Cell{Value: v, IsPreFilled: hasStartingValue},
				group: group[Loc{X: x, Y: y}],
			}
		}
	}

	g.Symbols = symbols
	if len(g.Symbols) == 0 {
		// Extract unique symbols and ensure they match the group count
		g.Symbols = make([]string, 0, len(symbolMap))
		for sym := range symbolMap {
			g.Symbols = append(g.Symbols, sym)
		}
	}
	if len(g.Symbols) == 0 {
		return fmt.Errorf("no symbols found in the provided cells")
	}

	groupVals := map[int]struct{}{}
	for _, g := range group {
		groupVals[g] = struct{}{}
	}
	if len(g.Symbols) != len(groupVals) {
		return fmt.Errorf("number of symbols (%d) does not match number of group values (%d)", len(g.Symbols), len(groupVals))
	}
	slices.Sort(g.Symbols)

	// Initialize options for empty cells
	for y := range g.Board {
		for x := range g.Board[y] {
			if g.Board[y][x].Cell.Value == "" {
				g.Board[y][x].Cell.Candidates = slices.Clone(g.Symbols)
			}
		}
	}

	err := g.EliminateCandidatesInit() // Initial elimination of candidates
	if err != nil {
		return fmt.Errorf("failed to initialize candidates: %w", err)
	}
	return nil
}

func (g *Game) FillBasic(cells [][]int) error {
	return g.FillInts(cells, DefaultGroup9x9, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"})
}
func (g *Game) FillBoard(board [][]GroupedCell) error {
	// Extract cells and groups from the provided board
	cells := make([][]string, len(board))
	group := make(map[Loc]int)
	symbolMap := map[string]struct{}{}

	for y := range board {
		cells[y] = make([]string, len(board[y]))
		for x := range board[y] {
			cells[y][x] = board[y][x].Cell.Value
			group[Loc{X: x, Y: y}] = board[y][x].group
			if board[y][x].Cell.Value != "" {
				symbolMap[board[y][x].Cell.Value] = struct{}{}
			}
		}
	}

	// Extract symbols
	symbols := make([]string, 0, len(symbolMap))
	for sym := range symbolMap {
		symbols = append(symbols, sym)
	}

	return g.Fill(cells, group, symbols)
}

func (g *Game) FillInts(cells [][]int, group map[Loc]int, symbols []string) error {
	strCells := make([][]string, len(cells))
	for y, row := range cells {
		strCells[y] = make([]string, len(row))
		for x, v := range row {
			if v == 0 {
				continue
			}
			strCells[y][x] = strconv.Itoa(v)
		}
	}
	return g.Fill(strCells, group, nil)
}

func (c *Cell) Set(v string) {
	c.Value = v
	c.Candidates = nil // Clear options since the cell is now filled
}

func (c *Cell) RemoveCandiates(vs []string) (removed []string) {
	if c.Value != "" {
		return nil // Cell is already filled, nothing to remove
	}

	// Remove empty vs values
	vs = slices.DeleteFunc(vs, func(v string) bool {
		return v == ""
	})

	removed = []string{}
	c.Candidates = slices.DeleteFunc(c.Candidates, func(c string) bool {
		if slices.Contains(vs, c) {
			removed = append(removed, c)
			return true
		}
		return false
	})

	return removed
}

func (g *Game) SingleCadidate() (x, y int, v string, ok bool) {
	for y := range g.Board {
		for x := range g.Board[y] {
			cell := g.Board[y][x].Cell
			if len(cell.Candidates) == 1 {
				return x, y, cell.Candidates[0], true
			}
		}
	}
	return 0, 0, "", false
}

func (g *Game) SetLastFilled(x, y int) {
	for iy := range g.Board {
		for jx := range g.Board[iy] {
			cell := g.Board[iy][jx].Cell
			if iy == y && jx == x {
				cell.IsLastFilled = true
				continue
			}
			cell.IsLastFilled = false
		}
	}
}

func (g *Game) Won() bool {
	for y := range g.Board {
		for x := range g.Board[y] {
			cell := g.Board[y][x].Cell
			if cell.Value == "" {
				return false // Found an empty cell
			}
		}
	}
	return true // All cells are filled and have no candidates
}

func (g *Game) BadBoard() error {
	rows, cols, groups := g.GetSectionedCells()
	// Check all sections (rows, columns, groups)
	sections := []struct {
		name  string
		cells [][]LocCell
	}{
		{"row", rows},
		{"column", cols},
		{"group", groups},
	}

	for _, section := range sections {
		for i, cells := range section.cells {
			values := make(map[string][]Loc)
			singleCandidates := make(map[string][]Loc)

			for _, lc := range cells {
				if lc.Cell.Value != "" {
					values[lc.Cell.Value] = append(values[lc.Cell.Value], lc.Loc)
				} else if len(lc.Cell.Candidates) == 1 {
					singleCandidates[lc.Cell.Candidates[0]] = append(singleCandidates[lc.Cell.Candidates[0]], lc.Loc)
				}
			}

			// Check for duplicate values
			for v, locs := range values {
				if len(locs) > 1 {
					return fmt.Errorf("duplicate value '%s' in %s %d at positions %v", v, section.name, i, locs)
				}
			}

			// Check for duplicate single candidates
			for v, locs := range singleCandidates {
				if len(locs) > 1 {
					return fmt.Errorf("multiple cells with only candidate '%s' in %s %d at positions %v", v, section.name, i, locs)
				}
			}
		}
	}

	return nil // Board is valid
}

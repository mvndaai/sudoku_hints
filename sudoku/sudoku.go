package sudoku

import (
	"fmt"
	"slices"
	"strconv"
)

type (
	GroupedCell struct {
		group int
		Cell  *Cell
	}

	LocCell struct {
		Loc  Loc
		Cell *Cell
	}

	Cell struct {
		value      string
		Candidates []string
	}

	Game struct {
		Symbols           []string
		Board             [][]GroupedCell
		HideSimple        bool
		RandomEliminators bool // If true, the eliminators will be run in a random order
	}
)

func (g *Game) Fill(cells [][]string, group map[Loc]int) error {
	symbols := map[string]struct{}{}

	// Fill in the cell values and track unique symbols
	g.Board = make([][]GroupedCell, len(cells))
	for y, row := range cells {
		g.Board[y] = make([]GroupedCell, len(row))
		for x, v := range row {
			if v != "" {
				symbols[v] = struct{}{}
			}
			g.Board[y][x] = GroupedCell{
				Cell:  &Cell{value: v},
				group: group[Loc{X: x, Y: y}],
			}
		}
	}

	// Extract unique symbols and ensure they match the group count
	g.Symbols = make([]string, 0, len(symbols))
	for sym := range symbols {
		g.Symbols = append(g.Symbols, sym)
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
			if g.Board[y][x].Cell.value == "" {
				g.Board[y][x].Cell.Candidates = slices.Clone(g.Symbols)
			}
		}
	}

	return nil
}

func (g *Game) FillBasic(cells [][]int) error {
	strCells := make([][]string, len(cells))
	for y, row := range cells {
		strCells[y] = make([]string, len(row))
		for x, v := range row {
			if v == 0 {
				continue
			}
			if v > 9 {
				return fmt.Errorf("invalid cell[x:%d,y:%d] value: %d, expected between 1 and 9", x, y, v)
			}
			strCells[y][x] = strconv.Itoa(v)
		}
	}
	return g.Fill(strCells, DefaultGropu9x9)
}

func (c *Cell) Set(v string) {
	c.value = v
	c.Candidates = nil // Clear options since the cell is now filled
}

func (c *Cell) RemoveCandiates(vs []string) (removed []string) {
	if c.value != "" {
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

func (g *Game) Won() bool {
	for y := range g.Board {
		for x := range g.Board[y] {
			cell := g.Board[y][x].Cell
			if cell.value == "" {
				return false // Found an empty cell
			}
		}
	}
	return true // All cells are filled and have no candidates
}

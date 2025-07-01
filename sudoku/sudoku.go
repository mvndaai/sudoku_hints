package sudoku

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/fatih/color"
)

type (
	GroupedCell struct {
		group int
		Cell  Cell
	}

	Cell struct {
		value      string
		candidates []string
	}

	Game struct {
		Symbols []string
		Board   [][]GroupedCell
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
				Cell:  Cell{value: v},
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
	if len(g.Symbols) != len(group) {
		return fmt.Errorf("number of symbols (%d) does not match number of groups (%d)", len(g.Symbols), len(group))
	}
	slices.Sort(g.Symbols)

	// Initialize options for empty cells
	for y := range g.Board {
		for x := range g.Board[y] {
			if g.Board[y][x].Cell.value == "" {
				g.Board[y][x].Cell.candidates = slices.Clone(g.Symbols)
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
	g.Fill(strCells, DefaultGropu9x9)
	return nil
}

func (g *Game) String() string {
	colors := []color.Attribute{color.FgYellow, color.FgBlue}

	var r string
	for y, row := range g.Board {
		if y != 0 {
			r += "\n"
		}
		for _, gc := range row {
			c := color.New(colors[gc.group%2])

			v := gc.Cell.value
			if v == "" {
				v = "⛝" // Use underscore for empty cells
			}
			r += c.Sprint(v) + " "
		}
	}
	return r
}

func (c *Cell) Set(v string) {
	c.value = v
	c.candidates = nil // Clear options since the cell is now filled
}

func (c *Cell) RemoveCandiates(vs []string) {
	// Remove empty vs values
	vs = slices.DeleteFunc(vs, func(v string) bool {
		return v == ""
	})

	c.candidates = slices.DeleteFunc(c.candidates, func(c string) bool {
		return slices.Contains(vs, c)
	})
}

func (g *Game) SingleCadidate() (x, y int, v string, ok bool) {
	for y := range g.Board {
		for x := range g.Board[y] {
			cell := &g.Board[y][x].Cell
			if len(cell.candidates) == 1 {
				return x, y, cell.candidates[0], true
			}
		}
	}
	return 0, 0, "", false
}

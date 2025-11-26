package sudoku

import (
	"fmt"
	"io"
	"log"

	"github.com/fatih/color"
)

// These are functions that print in the bash console

//github.com/gosuri/uilive

func (g *Game) String(lastUpdated *Loc) string {
	colors := []color.Attribute{color.FgYellow, color.FgCyan}

	var r string
	for y, row := range g.Board {
		if y != 0 {
			r += "\n"
		}
		for x, gc := range row {
			c := color.New(colors[gc.group%2])
			if lastUpdated != nil && lastUpdated.X == x && lastUpdated.Y == y {
				c = color.New(color.FgGreen, color.Bold) // Highlight the last updated cell
			}
			if !gc.Cell.IsPreFilled && gc.Cell.Value != "" {
				c.Add(color.Bold) // Make non-starting values bold
			}

			v := gc.Cell.Value
			if v == "" {
				v = "⛝" // Use underscore for empty cells
			}
			r += c.Sprint(v) + " "
		}
	}
	return r
}

//func (g *Game) StringWithCandidates() string {
//	colors := []color.Attribute{color.FgGreen, color.FgCyan}

//	var r string
//	for y, row := range g.Board {
//		if y != 0 {
//			r += "\n"
//		}
//		for _, gc := range row {
//			c := color.New(colors[gc.group%2])

//			v := gc.Cell.value
//			if v == "" {
//				v = "⛝" // Use underscore for empty cells
//			}
//			r += c.Sprint(v) + " "
//		}
//	}
//	return r
//}

var allChanges string

type gameWriter interface {
	io.Writer
	Flush() error
}

type scanner interface {
	Text() string
	Scan() bool
}

func (g *Game) StepThrough(w gameWriter, sc scanner) {
	log.Println("Starting StepThrough...", g.CellsWithRecentCandidates())
	var lastUpdated *Loc
	// Clear recent candidates from previous step
	if g.RunSimpleAfter {
		g.RemoveAllRecentCandidates()
	}
	if g.RunSimpleFirst {
		// Run simple eliminators first quietly
		_ = g.RemoveAllSimple(false)
		log.Println("After RunSimpleFirst RemoveAllSimple...", g.CellsWithRecentCandidates())
	}

	fmt.Fprintln(w, g.String(lastUpdated))
	solve := g.AutoSolve
	for {
		x, y, v, ok := g.SingleCadidate()
		if !ok {
			change, err := g.EliminateCandidates(false)
			if err != nil {
				w.Flush()
				fmt.Fprintln(w, g.String(lastUpdated))
				fmt.Fprint(w, color.New(color.FgRed).Sprintf("\nError: %v\n", err))
				//fmt.Fprint(w, allChanges)
				break
			}
			if change != "" {
				fmt.Fprintln(w, "Change: "+change)
				allChanges += change + "\n"
			}
			continue
		}

		f := fmt.Sprintf("Found single candidate at (x:%d, y:%d): %s\n", x, y, v)
		lastUpdated = &Loc{X: x, Y: y}
		fmt.Fprint(w, f)
		allChanges += f
		g.Board[y][x].Cell.Set(v)
		g.SetLastFilled(x, y)

		if err := g.BadBoard(); err != nil {
			w.Flush()
			fmt.Fprintln(w, g.String(lastUpdated))
			fmt.Fprint(w, color.New(color.FgRed).Sprintf("\nError: %v\n", err))
			fmt.Fprint(w, allChanges)
			break
		}

		if g.Won() {
			w.Flush()
			fmt.Fprintln(w, g.String(lastUpdated))
			fmt.Fprintln(w, color.New(color.FgGreen).Sprint("Congratulations! We solved the Sudoku puzzle!"))
			break
		}

		if !solve {
			if g.RunOnce {
				break
			}

			fmt.Fprint(w, color.New(color.FgYellow).Sprint("Enter to continue "))
			sc.Scan()
			if t := sc.Text(); t == "solve" || t == "s" {
				solve = true
			}
			fmt.Printf("\033[1A\033[K") // Move cursor up and clear the line - this is needed to avoid the console being cluttered with "Enter to continue" messages
			//fmt.Print("\033[H\033[2J") // Clear the console - needed because enter was breaking things
		}
		w.Flush()
		fmt.Fprintln(w, g.String(lastUpdated))
	}

	// After setting a value, remove candidates from other cells in the same row/col/group
	if g.RunSimpleAfter {
		log.Println("Calling RemoveAllSimple after setting cell...", g.CellsWithRecentCandidates())
		err := g.RemoveAllSimple(false)
		if err != nil {
			log.Println("Error in RemoveAllSimple:", err)
		}
		log.Println("RemoveAllSimple completed", g.CellsWithRecentCandidates())
	}
	log.Println("Finished StepThrough.", g.CellsWithRecentCandidates(), "==============================")
}

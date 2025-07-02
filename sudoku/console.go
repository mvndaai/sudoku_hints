package sudoku

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
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
			if !gc.Cell.startingValue && gc.Cell.value != "" {
				c.Add(color.Bold) // Make non-starting values bold
			}

			v := gc.Cell.value
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

func (g *Game) StepThroughConsole() {
	writer := uilive.New()
	writer.Start()
	scanner := bufio.NewScanner(os.Stdin)
	g.StepThrough(writer, scanner)
	writer.Stop()
}

func (g *Game) StepThrough(w gameWriter, sc scanner) {
	var lastUpdated *Loc
	if g.RunSimpleFirst {
		// Run simple eliminators first quietly
		for {
			change, err := g.EliminateCandidates(true)
			if err != nil {
				w.Flush()
				fmt.Fprintln(w, g.String(lastUpdated))
				fmt.Fprint(w, color.New(color.FgRed).Sprintf("\nError: %v\n", err))
				//fmt.Fprint(writer, allChanges)
				return
			}
			if change != "" {
				break
			}
		}
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

		f := fmt.Sprintf("Found single candidate at (%d, %d): %s\n", x, y, v)
		lastUpdated = &Loc{X: x, Y: y}
		fmt.Fprint(w, f)
		allChanges += f
		g.Board[y][x].Cell.Set(v)

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

}

package sudoku

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

// These are functions that print in the bash console

//github.com/gosuri/uilive

func (g *Game) String() string {
	colors := []color.Attribute{color.FgYellow, color.FgCyan}

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

func (g *Game) StepThrough() {
	writer := uilive.New()
	// start listening for updates and render
	writer.Start()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Fprintln(writer, g.String())
	var solve bool
	for {
		x, y, v, ok := g.SingleCadidate()
		if !ok {
			change, err := g.EliminateCandidates(false)
			if err != nil {
				writer.Flush()
				fmt.Fprintln(writer, g.String())
				fmt.Fprint(writer, color.New(color.FgRed).Sprintf("\nError: %v\n", err))
				//fmt.Fprint(writer, allChanges)
				break
			}
			if change != "" {
				fmt.Fprintln(writer, "Change: "+change)
				allChanges += change + "\n"
			}
			continue
		}

		f := fmt.Sprintf("Found single candidate at (%d, %d): %s\n", x, y, v)
		fmt.Fprint(writer, f)
		allChanges += f
		g.Board[y][x].Cell.Set(v)

		if err := g.BadBoard(); err != nil {
			writer.Flush()
			fmt.Fprintln(writer, g.String())
			fmt.Fprint(writer, color.New(color.FgRed).Sprintf("\nError: %v\n", err))
			fmt.Fprint(writer, allChanges)
			break
		}

		if g.Won() {
			writer.Flush()
			fmt.Fprintln(writer, g.String())
			fmt.Fprintln(writer, color.New(color.FgGreen).Sprint("Congratulations! We solved the Sudoku puzzle!"))
			break
		}

		if !solve {
			fmt.Fprint(writer, color.New(color.FgYellow).Sprint("Enter to continue "))
			scanner.Scan()
			if t := scanner.Text(); t == "solve" || t == "s" {
				solve = true
			}
			fmt.Printf("\033[1A\033[K") // Move cursor up and clear the line - this is needed to avoid the console being cluttered with "Enter to continue" messages
			//fmt.Print("\033[H\033[2J") // Clear the console - needed because enter was breaking things
		}
		writer.Flush()
		fmt.Fprintln(writer, g.String())
	}
	writer.Stop()
}

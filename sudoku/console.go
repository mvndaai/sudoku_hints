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

func (g *Game) StringWithCandidates() string {
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

func (g *Game) StepThrough() {
	writer := uilive.New()
	// start listening for updates and render
	writer.Start()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Fprintln(writer, g.String())
	for {
		x, y, v, ok := g.SingleCadidate()
		if !ok {
			change, err := g.EliminateCandidates()
			if err != nil {
				fmt.Fprint(writer, color.New(color.FgRed).Sprintf("\nError: %v\n", err))
				fmt.Fprintln(writer, g.String())
				break
			}
			fmt.Fprintln(writer, "Change: "+change)
			continue
		}

		fmt.Fprintf(writer, "Single candidate found at (%d, %d): %s\n", x, y, v)
		g.Board[y][x].Cell.Set(v)

		if g.Won() {
			fmt.Fprintln(writer, color.New(color.FgGreen).Sprint("Congratulations! We solved the Sudoku puzzle!"))
			fmt.Fprintln(writer, g.String())
			break
		}

		fmt.Fprintln(writer, color.New(color.FgYellow).Sprint("Enter to continue"))
		scanner.Scan()

		fmt.Fprintln(writer, g.String())
		writer.Flush()
	}

	//for i := 0; i <= 100; i++ {
	//	fmt.Fprintf(writer, "Downloading.. (%d/%d) GB\n", i, 100)
	//	time.Sleep(time.Millisecond * 5)
	//}

	//fmt.Fprintln(writer, "Finished: Downloaded 100GB")
	writer.Stop() // flush and stop rendering
}

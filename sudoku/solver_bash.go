//go:build !js && !wasm

package sudoku

import (
	"bufio"
	"os"

	"github.com/gosuri/uilive"
)

func (g *Game) StepThroughConsole() {
	writer := uilive.New()
	writer.Start()
	scanner := bufio.NewScanner(os.Stdin)
	g.StepThrough(writer, scanner)
	writer.Stop()
}

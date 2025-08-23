//go:build js && wasm
// +build js,wasm

package sudoku

import "io"

type writeFlusher struct {
	io.Writer
}

func (w *writeFlusher) Flush() error { return nil }

type noOpScanner struct{}

func (n *noOpScanner) Scan() bool   { return false }
func (n *noOpScanner) Text() string { return "" }

func (g *Game) StepThroughJavascript(w io.Writer) {
	writer := &writeFlusher{w}
	scanner := &noOpScanner{}
	g.StepThrough(writer, scanner)
}

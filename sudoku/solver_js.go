//go:build js && wasm
// +build js,wasm

package sudoku

import "io"

type writeFlusher struct {
	io.Writer
}

func (w *writeFlusher) Flush() error { return nil }

type noOpWriter struct{}

func (n *noOpWriter) Write(p []byte) (int, error) { return 0, nil }

type noOpScanner struct{}

func (n *noOpScanner) Scan() bool   { return false }
func (n *noOpScanner) Text() string { return "" }

func (g *Game) StepThroughJavascript(w io.Writer) {
	if w == nil {
		w = &noOpWriter{}
	}
	writer := &writeFlusher{w}
	scanner := &noOpScanner{}
	g.StepThrough(writer, scanner)
}

//go:build js && wasm
// +build js,wasm

package sudoku

import (
	"fmt"
	"io"
)

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
	g.RunOnce = true
	g.StepThrough(writer, scanner)
}

func Next(g *Game, eliminators []string) error {
	g.RemoveAllRecentCandidates()

	err := g.RemoveAllSimple(false)
	if err != nil {
		return fmt.Errorf("failed to remove all simple candidates: %w", err)
	}

	g.StepThroughJavascript(nil)
	return nil
}

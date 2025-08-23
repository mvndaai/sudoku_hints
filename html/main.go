//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/mvndaai/sudoku_hints/sudoku"
	"github.com/mvndaai/sudoku_hints/sudoku/boards"
)

func main() {
	var m = make(map[string]any)
	m["getKey"] = getKey()
	m["random"] = getRandomBoard()
	m["convertOCR"] = convertOCR()
	m["next"] = next()

	js.Global().Set("golang", m)

	// Keep the program alive so functions can be run over and over
	<-make(chan bool)
}

// RapidAPIKey is set via build flag: -ldflags "-X 'main.RapidAPIKey=YOUR_KEY'"
var RapidAPIKey string

func getKey() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		return RapidAPIKey
	})
}

func getRandomBoard() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		g := sudoku.Game{}
		err := g.FillBasic(boards.RandomBasicBoard())
		if err != nil {
			return err
		}

		b, err := json.Marshal(g.Board)
		if err != nil {
			return err
		}
		return string(b)
	})
}

func convertOCR() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {

		board, err := sudoku.ConvertFromOCRFormat(args[0].String())
		if err != nil {
			return err
		}

		g := sudoku.Game{}
		g.FillBasic(board)
		b, err := json.Marshal(g.Board)
		if err != nil {
			return err
		}
		return string(b)
	})
}

func next() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		board := [][]sudoku.GroupedCell{}
		err := json.Unmarshal([]byte(args[0].String()), &board)
		if err != nil {
			return err
		}

		g := sudoku.Game{}
		g.FillBoard(board)
		g.RunOnce = true

		g.StepThroughJavascript(nil)

		b, err := json.Marshal(g.Board)
		if err != nil {
			return err
		}
		return string(b)
	})
}

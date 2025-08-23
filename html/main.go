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

		//board := boards.RandomBasicBoard()

		err := g.FillBasic(boards.RandomBasicBoard())
		if err != nil {
			return err
		}

		b, err := json.Marshal(g.Board)
		if err != nil {
			return err
		}

		return string(b)

		//return board
	})
}

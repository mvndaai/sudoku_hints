//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
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
	m["processOCR"] = processOCR()

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

func processOCR() js.Func { // If you have an http request it needs to return a promise
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 2 {
			return js.Global().Get("Promise").Call("reject", "Invalid number of arguments passed")
		}

		// Capture the original arguments from processOCR
		filename := args[0].String()
		jsBytes := args[1]

		handler := js.FuncOf(func(this js.Value, args []js.Value) any {
			resolve := args[0]
			reject := args[1]

			go func() {
				// Use the captured arguments from the outer scope
				length := jsBytes.Get("length").Int()
				fileBytes := make([]byte, length)
				js.CopyBytesToGo(fileBytes, jsBytes)

				resp, err := sudoku.ProcessImage(RapidAPIKey, filename, fileBytes)
				if err != nil {
					reject.Invoke(err.Error())
					return
				}

				pretty, err := json.Marshal(resp)
				if err != nil {
					reject.Invoke(fmt.Errorf("Could not marshal json: %w %v", err, resp))
					return
				}
				resolve.Invoke(string(pretty))
			}()

			return nil
		})

		return js.Global().Get("Promise").New(handler)
	})
}

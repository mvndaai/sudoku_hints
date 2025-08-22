//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"
)

func main() {
	var m = make(map[string]any)
	m["getKey"] = getKey()

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

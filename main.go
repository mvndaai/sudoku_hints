package main

import (
	"log"

	"github.com/mvndaai/sudoku_hints/sudoku"
	"github.com/mvndaai/sudoku_hints/sudoku/boards"
)

func main() {
	g := sudoku.Game{}
	err := g.FillBasic(boards.NYTHard2June2025)
	//err := g.FillBasic(basicEasy)
	if err != nil {
		log.Fatalf("Failed to fill game: %v", err)
	}

	//g.HideSimple = true // Hide basic eliminators
	//g.RandomEliminators = true // Randomize the order of eliminators
	g.RunSimpleFirst = true // Run simple eliminators first quietly
	g.StepThrough()
}

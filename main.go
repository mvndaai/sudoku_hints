package main

import (
	"log"

	"github.com/mvndaai/sudoku_hints/sudoku"
	"github.com/mvndaai/sudoku_hints/sudoku/boards"
)

func main() {
	g := sudoku.Game{}
	err := g.FillBasic(boards.NYTHard2June2025)
	//err := g.FillBasic(boards.BasicEasy)
	//err := g.FillBasic(boards.BasicHard)
	if err != nil {
		log.Fatalf("Failed to fill game: %v", err)
	}

	//g.HideSimple = true // Hide basic eliminators
	//g.RandomEliminators = true // Randomize the order of eliminators TODO this causes errors.
	g.RunSimpleFirst = true // Run simple eliminators first quietly
	g.AutoSolve = true      // Automatically solve the game
	g.StepThroughConsole()
}

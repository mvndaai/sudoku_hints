package main

import (
	"log"

	"github.com/mvndaai/sudoku_hints/sudoku"
)

var (
	basicEasy = [][]int{
		{8, 0, 1, 0, 3, 5, 0, 2, 0},
		{0, 0, 0, 2, 7, 6, 0, 5, 1},
		{0, 0, 6, 9, 0, 1, 0, 7, 3},

		{0, 9, 8, 0, 1, 0, 0, 3, 4},
		{7, 6, 0, 3, 5, 0, 0, 0, 0},
		{1, 0, 0, 0, 4, 9, 6, 0, 0},

		{0, 0, 0, 0, 9, 0, 5, 0, 0},
		{0, 1, 0, 0, 6, 0, 0, 0, 0},
		{6, 8, 3, 5, 0, 0, 1, 9, 0},
	}

	basicHard = [][]int{
		{5, 0, 0, 0, 2, 7, 0, 0, 0},
		{3, 0, 0, 0, 0, 0, 5, 0, 6},
		{0, 4, 0, 3, 0, 0, 0, 0, 0},

		{6, 9, 0, 0, 0, 2, 0, 0, 0},
		{0, 0, 1, 0, 9, 0, 0, 0, 0},
		{0, 0, 0, 8, 0, 0, 0, 0, 5},

		{0, 0, 8, 0, 0, 0, 0, 9, 0},
		{4, 0, 0, 0, 0, 6, 0, 0, 1},
		{0, 0, 0, 0, 0, 1, 0, 7, 0},
	}
)

func main() {
	_ = basicEasy // Use this to avoid unused variable error
	_ = basicHard // Use this to avoid unused variable error

	g := sudoku.Game{}
	err := g.FillBasic(basicHard)
	//err := g.FillBasic(basicEasy)
	if err != nil {
		log.Fatalf("Failed to fill game: %v", err)
	}

	g.HideSimple = true        // Hide basic eliminators
	g.RandomEliminators = true // Randomize the order of eliminators
	g.StepThrough()
}

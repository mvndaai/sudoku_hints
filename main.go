package main

import (
	"fmt"
	"log"

	"github.com/mvndaai/sudoku_hints/sudoku"
)

func main() {
	log.Print("Hello, World!")

	g := sudoku.Game{}
	g.FillBasic([][]int{
		{8, 0, 1, 0, 3, 5, 0, 2, 0},
		{0, 0, 0, 2, 7, 6, 0, 5, 1},
		{0, 0, 6, 9, 0, 1, 0, 7, 3},

		{0, 9, 8, 0, 1, 0, 0, 3, 4},
		{7, 6, 0, 3, 5, 0, 0, 0, 0},
		{1, 0, 0, 0, 4, 9, 6, 0, 0},

		{0, 0, 0, 0, 9, 0, 5, 0, 0},
		{0, 1, 0, 0, 6, 0, 0, 0, 0},
		{6, 8, 3, 5, 0, 0, 1, 9, 0},
	})
	g.Board[0][8].Cell.Set("6")
	fmt.Println(g.String())
}

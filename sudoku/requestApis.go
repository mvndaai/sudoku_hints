package sudoku

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type dosukoResponse struct {
	Newboard struct {
		Grids []struct {
			Value      [][]int `json:"value"`
			Difficulty string  `json:"difficulty"`
		} `json:"grids"`
	} `json:"newboard"`
}

func RequestDosukoPuzzle() (Game, error) {
	resp, err := http.Get("https://sudoku-api.vercel.app/api/dosuku?query={newboard(limit:1){grids{value,difficulty}}}")
	if err != nil {
		return Game{}, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return Game{}, fmt.Errorf("failed to get puzzle: %s", resp.Status)
	}

	var puzzleResponse dosukoResponse
	if err := json.NewDecoder(resp.Body).Decode(&puzzleResponse); err != nil {
		return Game{}, err
	}

	g := Game{}
	g.FillBasic(puzzleResponse.Newboard.Grids[0].Value)
	g.Difficulty = puzzleResponse.Newboard.Grids[0].Difficulty
	return g, nil
}

// https://www.youdosudoku.com/
// https://www.api-ninjas.com/api/sudoku
// https://sudoku-game-and-api.netlify.app/api/sudoku

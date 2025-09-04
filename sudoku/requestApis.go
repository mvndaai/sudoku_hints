package sudoku

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func RequestYouDoSudokuPuzzle(difficulty string) (Game, error) {
	// Create request body
	requestBody := map[string]interface{}{
		"difficulty": difficulty,
		"solution":   false,
		"array":      true,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return Game{}, err
	}

	resp, err := http.Post("https://youdosudoku.com/api/", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return Game{}, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return Game{}, fmt.Errorf("failed to get puzzle: %s", resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Game{}, err
	}

	// Define response structure
	type Response struct {
		Difficulty string  `json:"difficulty"`
		Puzzle     [][]int `json:"puzzle"`
		//Solution   string `json:"solution"`
	}

	// Parse response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Game{}, err
	}

	//// Convert puzzle string to [9][9]int
	//var puzzleArray [9][9]int
	//for i := 0; i < 81; i++ {
	//	row := i / 9
	//	col := i % 9
	//	digit := response.Puzzle[i] - '0'
	//	puzzleArray[row][col] = int(digit)
	//}

	g := Game{}
	g.FillBasic(response.Puzzle)
	g.Difficulty = response.Difficulty
	return g, nil
}

// https://www.api-ninjas.com/api/sudoku
// https://sudoku-game-and-api.netlify.app/api/sudoku

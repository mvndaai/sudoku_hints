package sudoku

import (
	"encoding/json"
	"fmt"
)

type ocrAPI struct {
	Puzzle struct {
		Rows []struct {
			Cells []struct {
				CellType   string `json:"cell_type"`
				Value      int    `json:"value,omitempty"`
				Candidates []int  `json:"candidates,omitempty"`
			} `json:"cells"`
		} `json:"rows"`
	} `json:"puzzle"`
}

func ConvertFromOCRFormat(s string) ([][]int, error) {
	oa := ocrAPI{}
	err := json.Unmarshal([]byte(s), &oa)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall ocr api format %w", err)
	}

	r := [][]int{}
	for _, row := range oa.Puzzle.Rows {
		rv := []int{}
		for _, cell := range row.Cells {
			rv = append(rv, cell.Value)
		}
		r = append(r, rv)
	}

	return r, nil
}

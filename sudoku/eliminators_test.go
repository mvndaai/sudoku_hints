package sudoku

import (
	"encoding/json"
	"log"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEliminatorFilledCell(t *testing.T) {
	cells := []LocCell{
		{Loc: Loc{X: 0, Y: 0}, Cell: &Cell{Value: "1", Candidates: nil}},
		{Loc: Loc{X: 1, Y: 0}, Cell: &Cell{Value: "2", Candidates: nil}},
		{Loc: Loc{X: 2, Y: 0}, Cell: &Cell{Value: "3", Candidates: nil}},
		{Loc: Loc{X: 3, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}},
		{Loc: Loc{X: 4, Y: 0}, Cell: &Cell{Value: "5", Candidates: nil}},
		{Loc: Loc{X: 5, Y: 0}, Cell: &Cell{Value: "6", Candidates: nil}},
		{Loc: Loc{X: 6, Y: 0}, Cell: &Cell{Value: "7", Candidates: nil}},
		{Loc: Loc{X: 7, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}},
		{Loc: Loc{X: 8, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}},
	}

	expectedChanges := []string{
		"removed candidates (x:3,y:0) [1 2 3 5 6 7]",
		"removed candidates (x:7,y:0) [1 2 3 5 6 7]",
		"removed candidates (x:8,y:0) [1 2 3 5 6 7]",
		"",
	}

	for _, expectedChange := range expectedChanges {
		change, err := EliminatorFilledCell.PartitionEliminator(cells)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if change != expectedChange {
			t.Errorf("expected change '%s', got '%s'", expectedChange, change)
		}
	}
}

func TestEliminatorEliminatorUniqueCandidate(t *testing.T) {
	cells := []LocCell{
		{Loc: Loc{X: 0, Y: 0}, Cell: &Cell{Value: "1", Candidates: nil}},
		{Loc: Loc{X: 1, Y: 0}, Cell: &Cell{Value: "2", Candidates: nil}},
		{Loc: Loc{X: 2, Y: 0}, Cell: &Cell{Value: "3", Candidates: nil}},
		{Loc: Loc{X: 3, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"4", "5"}}},
		{Loc: Loc{X: 4, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"4", "5"}}},
		{Loc: Loc{X: 5, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"4", "5", "6"}}},
		{Loc: Loc{X: 6, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"7"}}},
		{Loc: Loc{X: 7, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"8", "9"}}},
		{Loc: Loc{X: 8, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"8", "9"}}},
	}

	expectedChanges := []string{ // Note these can come in any order
		"removed candidates (x:5,y:0) [4 5]",
	}

	foundChanges := []string{}
	for {
		b, _ := json.Marshal(cells)
		t.Logf("%s", b)
		change, err := EliminatorUniqueCandidate.PartitionEliminator(cells)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if change == "" {
			break
		}
		foundChanges = append(foundChanges, change)
		t.Log(change)
	}

	slices.Sort(foundChanges)
	slices.Sort(expectedChanges)
	assert.EqualValues(t, expectedChanges, foundChanges)
}

func TestEliminatorCandidateChains(t *testing.T) {
	cells := []LocCell{
		{Loc: Loc{X: 0, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"1", "2"}}},
		{Loc: Loc{X: 1, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"1", "2"}}},
		{Loc: Loc{X: 2, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"3", "4"}}},
		{Loc: Loc{X: 3, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"3", "4"}}},
		{Loc: Loc{X: 4, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"4", "5", "6", "7"}}},
		{Loc: Loc{X: 5, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"5", "6", "7"}}},
		{Loc: Loc{X: 6, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"5", "6", "7"}}},
		{Loc: Loc{X: 7, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"2", "3", "4", "5", "6", "7", "8", "9"}}},
		{Loc: Loc{X: 8, Y: 0}, Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}},
	}

	expectedChanges := []string{ // Note these can come in any order
		"removed candidates (x:4,y:0) [4] from chain of size 2",
		"removed candidates (x:7,y:0) [2] from chain of size 2",
		"removed candidates (x:7,y:0) [3 4] from chain of size 2",
		"removed candidates (x:7,y:0) [5 6 7] from chain of size 3",
		"removed candidates (x:8,y:0) [1 2] from chain of size 2",
		"removed candidates (x:8,y:0) [3 4] from chain of size 2",
		"removed candidates (x:8,y:0) [5 6 7] from chain of size 3",
	}

	foundChanges := []string{}
	for {
		b, _ := json.Marshal(cells)
		t.Logf("%s", b)
		change, err := EliminatorCandidateChains.PartitionEliminator(cells)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if change == "" {
			break
		}
		foundChanges = append(foundChanges, change)
		t.Log(change)
	}

	slices.Sort(foundChanges)
	slices.Sort(expectedChanges)
	if !slices.Equal(foundChanges, expectedChanges) {
		t.Errorf("changes\nexpected %#v\ngot      %#v", expectedChanges, foundChanges)
	}
}

func TestEliminatorGroupAndRowColumn(t *testing.T) {
	tests := []struct {
		name     string
		board    [][]int
		expected []string // Expected changes
	}{
		{
			name: "Rows",
			board: [][]int{
				{9, 8, 7, 1, 2, 3, 0, 0, 0},
				{6, 5, 4, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},

				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},

				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			expected: []string{
				"removed candidates (x:6,y:2) [1 2 3]",
				"removed candidates (x:7,y:2) [1 2 3]",
				"removed candidates (x:8,y:2) [1 2 3]",
			},
		},
		//{
		//	name: "Columns",
		//	board: [][]int{
		//		{9, 6, 0, 0, 0, 0, 0, 0, 0},
		//		{8, 5, 0, 0, 0, 0, 0, 0, 0},
		//		{7, 4, 0, 0, 0, 0, 0, 0, 0},

		//		{1, 0, 0, 0, 0, 0, 0, 0, 0},
		//		{2, 0, 0, 0, 0, 0, 0, 0, 0},
		//		{3, 0, 0, 0, 0, 0, 0, 0, 0},

		//		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		//	},
		//	expected: []string{
		//		"removed candidates (x:2,y:6) [1 2 3]",
		//		"removed candidates (x:2,y:7) [1 2 3]",
		//		"removed candidates (x:2,y:8) [1 2 3]",
		//	},
		//},
		//{
		//	name: "Hard",
		//	board: [][]int{
		//		{5, 0, 0, 0, 2, 7, 0, 0, 0},
		//		{3, 0, 0, 0, 0, 0, 5, 0, 6},
		//		{0, 4, 0, 3, 0, 0, 0, 0, 0},

		//		{6, 9, 0, 0, 0, 2, 0, 0, 0},
		//		{0, 0, 1, 0, 9, 0, 0, 0, 0},
		//		{0, 0, 0, 8, 0, 0, 0, 0, 5},

		//		{0, 0, 8, 0, 0, 0, 0, 9, 0},
		//		{4, 0, 0, 0, 0, 6, 0, 0, 1},
		//		{0, 0, 0, 0, 0, 1, 0, 7, 0},
		//	},
		//	expected: []string{
		//		"removed candidates (x:1,y:7) [5]",
		//		"removed candidates (x:2,y:7) [5]",
		//		"removed candidates (x:3,y:0) [9]",
		//		"removed candidates (x:3,y:1) [9]",
		//		"removed candidates (x:3,y:7) [5]",
		//		"removed candidates (x:4,y:1) [8]",
		//		"removed candidates (x:4,y:2) [8]",
		//		"removed candidates (x:4,y:7) [5]",
		//		"removed candidates (x:6,y:0) [9]",
		//		"removed candidates (x:6,y:2) [9]",
		//		"removed candidates (x:6,y:4) [6]",
		//		"removed candidates (x:6,y:4) [8]",
		//		"removed candidates (x:6,y:5) [6]",
		//		"removed candidates (x:7,y:4) [8]",
		//		"removed candidates (x:8,y:4) [8]",
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{}
			err := g.FillBasic(tt.board)
			require.NoError(t, err)
			for {
				if _, err := g.EliminateCandidates(true); err != nil {
					break
				}
			}

			foundChanges := []string{}
			for {
				change, err := EliminatorGroupAndRowColumn.GameEliminator(g)
				require.NoError(t, err)
				if change == "" {
					break
				}
				foundChanges = append(foundChanges, change)
				t.Log(change)
				log.Println(change)

				require.NoError(t, g.BadBoard())
			}

			slices.Sort(foundChanges)
			slices.Sort(tt.expected)
			assert.EqualValues(t, tt.expected, foundChanges)
		})
	}

}

func TestEliminatorFistemafelRing(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Game
		expected []string
	}{
		{
			name: "One group complete removes candidates from other",
			setup: func() *Game {
				g := &Game{}
				g.Symbols = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}

				// Create a board with corner group (0) complete with all 16 cells filled with unique values 1-8 (two of each)
				// and ring group (1) with 9 as a candidate that should be removed
				g.Board = [][]GroupedCell{
					{{Cell: &Cell{Value: "1", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "2", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "3", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "4", IsPreFilled: true}, group: 0}},
					{{Cell: &Cell{Value: "5", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "6", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "7", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "8", IsPreFilled: true}, group: 0}},
					{{Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}},
					{{Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}},
					{{Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}},
					{{Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}},
					{{Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 1}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}},
					{{Cell: &Cell{Value: "1", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "2", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "3", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "4", IsPreFilled: true}, group: 0}},
					{{Cell: &Cell{Value: "5", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "6", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}, group: 2}, {Cell: &Cell{Value: "7", IsPreFilled: true}, group: 0}, {Cell: &Cell{Value: "8", IsPreFilled: true}, group: 0}},
				}
				return g
			},
			expected: []string{
				"removed [9] from ring cell at (2,2) because corners is complete without these values",
				"removed [9] from ring cell at (3,2) because corners is complete without these values",
				"removed [9] from ring cell at (4,2) because corners is complete without these values",
				"removed [9] from ring cell at (5,2) because corners is complete without these values",
				"removed [9] from ring cell at (6,2) because corners is complete without these values",
				"removed [9] from ring cell at (2,3) because corners is complete without these values",
				"removed [9] from ring cell at (6,3) because corners is complete without these values",
				"removed [9] from ring cell at (2,4) because corners is complete without these values",
				"removed [9] from ring cell at (6,4) because corners is complete without these values",
				"removed [9] from ring cell at (2,5) because corners is complete without these values",
				"removed [9] from ring cell at (6,5) because corners is complete without these values",
				"removed [9] from ring cell at (2,6) because corners is complete without these values",
				"removed [9] from ring cell at (3,6) because corners is complete without these values",
				"removed [9] from ring cell at (4,6) because corners is complete without these values",
				"removed [9] from ring cell at (5,6) because corners is complete without these values",
				"removed [9] from ring cell at (6,6) because corners is complete without these values",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.setup()

			foundChanges := []string{}
			for {
				change, err := EliminatorFistemafelRing.GameEliminator(g)
				require.NoError(t, err)
				if change == "" {
					break
				}
				foundChanges = append(foundChanges, change)
				t.Log(change)
			}

			assert.ElementsMatch(t, tt.expected, foundChanges)
		})
	}
}

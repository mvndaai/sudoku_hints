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

package sudoku

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestEliminatorFilledCell(t *testing.T) {
	cells := []LocCell{
		{Loc: Loc{X: 0, Y: 0}, Cell: &Cell{value: "1", Candidates: nil}},
		{Loc: Loc{X: 1, Y: 0}, Cell: &Cell{value: "2", Candidates: nil}},
		{Loc: Loc{X: 2, Y: 0}, Cell: &Cell{value: "3", Candidates: nil}},
		{Loc: Loc{X: 3, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}},
		{Loc: Loc{X: 4, Y: 0}, Cell: &Cell{value: "5", Candidates: nil}},
		{Loc: Loc{X: 5, Y: 0}, Cell: &Cell{value: "6", Candidates: nil}},
		{Loc: Loc{X: 6, Y: 0}, Cell: &Cell{value: "7", Candidates: nil}},
		{Loc: Loc{X: 7, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}},
		{Loc: Loc{X: 8, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}},
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

func TestEliminatorMatchingCandidates(t *testing.T) {
	cells := []LocCell{
		{Loc: Loc{X: 0, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"1", "2"}}},
		{Loc: Loc{X: 1, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"1", "2"}}},
		{Loc: Loc{X: 2, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"3", "4"}}},
		{Loc: Loc{X: 3, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"3", "4"}}},
		{Loc: Loc{X: 4, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"4", "5", "6", "7"}}},
		{Loc: Loc{X: 5, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"5", "6", "7"}}},
		{Loc: Loc{X: 6, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"5", "6", "7"}}},
		{Loc: Loc{X: 7, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"2", "3", "4", "5", "6", "7", "8", "9"}}},
		{Loc: Loc{X: 8, Y: 0}, Cell: &Cell{value: "", Candidates: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}}},
	}

	expectedChanges := []string{ // Note these can come in any order
		"removed candidates (x:4,y:0) [4]",
		"removed candidates (x:7,y:0) [2]",
		"removed candidates (x:7,y:0) [3 4]",
		"removed candidates (x:7,y:0) [5 6 7]",
		"removed candidates (x:8,y:0) [1 2]",
		"removed candidates (x:8,y:0) [3 4]",
		"removed candidates (x:8,y:0) [5 6 7]",
	}

	foundChanges := []string{}
	for {
		b, _ := json.Marshal(cells)
		t.Logf("%s", b)
		change, err := EliminatorMatchingCandidates.PartitionEliminator(cells)
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

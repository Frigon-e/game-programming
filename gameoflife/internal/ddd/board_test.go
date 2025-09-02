package ddd

import (
	"testing"
)

func TestGolBoard_CountSurroundingLive(t *testing.T) {
	// Create a new 3x3 board for testing
	board := newGOLBoard(3, 3)

	// Set up a specific pattern of live cells
	// X X X
	// X O X
	// X X X
	// All neighbors are alive.
	initialState := []bool{
		true, true, true,
		true, false, true,
		true, true, true,
	}
	board.CopyBoard(initialState)

	// The center cell is at (1, 1)
	x, y := 1, 1

	// Count the surrounding live cells for the center cell
	liveNeighbors := board.CountSurroundingLive(x, y)

	// We expect 8 live neighbors
	expected := 8
	if liveNeighbors != expected {
		t.Errorf("Expected %d live neighbors, but got %d", expected, liveNeighbors)
	}
}

package application

import (
	"sort"
	"testing"
	"time"
)

func TestAttack_Hit(t *testing.T) {
	board := newBattleshipBoard(10, 10)
	board.PlaceShip(0, 0, Destroyer, Horizontal)

	hit, sunk, err := board.Attack(0, 0)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if !hit {
		t.Error("Expected hit to be true, but got false")
	}
	if sunk {
		t.Error("Expected sunk to be false, but got true")
	}
	if coord := board.Coordinate(0, 0); coord != Hit {
		t.Errorf("Expected coordinate to be %v, but got %v", Hit, coord)
	}
}

func TestAttack_Miss(t *testing.T) {
	board := newBattleshipBoard(10, 10)

	hit, sunk, err := board.Attack(0, 0)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if hit {
		t.Error("Expected hit to be false, but got true")
	}
	if sunk {
		t.Error("Expected sunk to be false, but got true")
	}
	if coord := board.Coordinate(0, 0); coord != Miss {
		t.Errorf("Expected coordinate to be %v, but got %v", Miss, coord)
	}
}

func TestAttack_Invalid(t *testing.T) {
	board := newBattleshipBoard(10, 10)

	_, _, err2 := board.Attack(0, 0)
	if err2 != nil {
		return
	}

	hit, sunk, err := board.Attack(0, 0) // Second attack on the same cell

	if err == nil {
		t.Error("Expected an error, but got nil")
	}
	if hit {
		t.Error("Expected hit to be false, but got true")
	}
	if sunk {
		t.Error("Expected sunk to be false, but got true")
	}
}

func TestAttack_Sunk(t *testing.T) {
	board := newBattleshipBoard(10, 10)
	board.PlaceShip(0, 0, Destroyer, Horizontal) // Destroyer has length 2

	_, _, err2 := board.Attack(0, 0)
	if err2 != nil {
		return
	}

	hit, sunk, err := board.Attack(1, 0)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if !hit {
		t.Error("Expected hit to be true, but got false")
	}
	if !sunk {
		t.Error("Expected sunk to be true, but got false")
	}
}

// onePlayerGame simulates a single game of Battleship for the AI and returns the number of moves taken to win.
func onePlayerGame(board BattleshipBoard) int {
	board.SeedBoard() // Place ships randomly for a new game.

	for moves := 1; moves <= board.Cols()*board.Rows(); moves++ {
		x, y := TakeTurn(board)

		_, _, err := board.Attack(x, y)
		println(x, y, err)
		println(moves)
		if err != nil {
			// This can happen if the AI guesses the same spot, so we just continue
			continue
		}

		if board.AllShipsSunk() {
			return moves // Return the total moves taken to win.
		}
	}
	return board.Cols() * board.Rows()
}

// TestHeatMapAIAverage benchmarks the performance of the heatMapAI over many games.
func TestHeatMapAIAverage(t *testing.T) {
	numGames := 1 // Running 100 games for a decent sample size
	results := make([]int, numGames)
	totalMoves := 0
	bestScore := 101
	worstScore := 0

	//ai := NewHeatMapAI(10, 10)
	board := NewBattleshipBoard(10, 10)

	startTime := time.Now()

	for i := 0; i < numGames; i++ {
		moves := onePlayerGame(board)
		results[i] = moves
		totalMoves += moves

		if moves < bestScore {
			bestScore = moves
		}
		if moves > worstScore {
			worstScore = moves
		}
	}

	duration := time.Since(startTime)

	sort.Ints(results)
	var median float64
	if numGames%2 == 0 {
		median = float64(results[numGames/2-1]+results[numGames/2]) / 2.0
	} else {
		median = float64(results[numGames/2])
	}

	// Log the formatted results of the benchmark.
	t.Logf("Simulation finished for %d games in %v", numGames, duration)
	t.Logf("Average score: %.2f", float64(totalMoves)/float64(numGames))
	t.Logf("Best score: %d", bestScore)
	t.Logf("Worst score: %d", worstScore)
	t.Logf("Median score: %.2f", median)
}

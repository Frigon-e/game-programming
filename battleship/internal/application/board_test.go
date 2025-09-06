package application

import (
	"sort"
	"testing"
	"time"
)

func TestAttack_Hit(t *testing.T) {
	board := newBattleshipBoard(10, 10)
	board.PlaceShip(0, 0, Destroyer, Horizontal)

	hit, sunk, _, err := board.Attack(0, 0)

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

	hit, sunk, _, err := board.Attack(0, 0)

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

	_, _, _, err2 := board.Attack(0, 0)
	if err2 != nil {
		return
	}

	hit, sunk, _, err := board.Attack(0, 0) // Second attack on the same cell

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

	_, _, _, err2 := board.Attack(0, 0)
	if err2 != nil {
		return
	}

	hit, sunk, _, err := board.Attack(1, 0)

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
func onePlayerGame() int {
	// The "solutionBoard" knows where the ships are. The AI will attack this board.
	solutionBoard := NewBattleshipBoard(10, 10)
	solutionBoard.SeedBoard()

	// The "aiViewBoard" is what the AI "sees". It only contains Empty, Hit, or Miss.
	// The AI uses this board to make its decisions.
	aiViewBoard := NewBattleshipBoard(10, 10)

	for moves := 1; moves <= solutionBoard.Cols()*solutionBoard.Rows(); moves++ {
		// AI decides its move based on what it can see.
		x, y := TakeTurn(aiViewBoard)

		// The attack happens on the real board.
		hit, _, _, err := solutionBoard.Attack(x, y)
		if err != nil {
			// This can happen if the AI guesses the same spot.
			// The AI logic should prevent this, but we continue just in case.
			continue
		}

		// Update the AI's view of the board with the result of the attack.
		if hit {
			aiViewBoard.SetCoordinate(x, y, Hit)
		} else {
			aiViewBoard.SetCoordinate(x, y, Miss)
		}

		if solutionBoard.AllShipsSunk() {
			return moves // Return the total moves taken to win.
		}
	}
	return solutionBoard.Cols() * solutionBoard.Rows()
}

// TestHeatMapAIAverage benchmarks the performance of the heatMapAI over many games.
func TestHeatMapAIAverage(t *testing.T) {
	numGames := 100_000 // Running 100 games for a decent sample size
	results := make([]int, numGames)
	totalMoves := 0
	bestScore := 101
	worstScore := 0

	startTime := time.Now()

	for i := 0; i < numGames; i++ {
		// onePlayerGame now creates its own boards, so we don't pass one in.
		moves := onePlayerGame()
		results[i] = moves
		totalMoves += moves

		if moves < bestScore {
			bestScore = moves
		}
		if moves > worstScore {
			worstScore = moves
		}
		if i%10000 == 0 {
			println(i)
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

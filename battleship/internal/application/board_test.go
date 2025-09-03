package application

import (
	"testing"
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
	board.Attack(0, 0) // First attack

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

	board.Attack(0, 0)
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

package application

import (
	"SideProjectGames/internal/ddd"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	Empty      uint8 = 0
	Hit        uint8 = 1
	Miss       uint8 = 2
	SUNK       uint8 = 3
	Carrier    uint8 = 4
	Battleship uint8 = 5
	Cruiser    uint8 = 6
	Submarine  uint8 = 7
	Destroyer  uint8 = 8
)

const (
	Horizontal uint8 = 0
	Vertical   uint8 = 1
)

type BattleshipBoard interface {
	ddd.Board[uint8]
	SeedBoard()
	Attack(x, y int) (hit, sunk bool, shipType uint8, err error)
	PlaceShip(x int, y int, shipType uint8, orientation uint8) bool
	IsCellSunk(x, y int) bool
	IsShipSunk(ship uint8) (sunk bool, err error)
	AllShipsSunk() bool
	RecordSunkShip(shipType uint8)
	HitShipAt() map[[2]int]uint8
	SunkShips() map[uint8]bool
	CopyHitValues(otherBoard BattleshipBoard)
}

type battleshipBoard struct {
	ddd.Board[uint8]
	sunkShips map[uint8]bool
	hitShipAt map[[2]int]uint8
}

var _ BattleshipBoard = (*battleshipBoard)(nil)

func NewBattleshipBoard(width int, height int) BattleshipBoard {
	return newBattleshipBoard(width, height)
}

func newBattleshipBoard(width int, height int) BattleshipBoard {
	return &battleshipBoard{Board: ddd.NewBoard[uint8](width, height), sunkShips: map[uint8]bool{
		Carrier:    false,
		Battleship: false,
		Cruiser:    false,
		Submarine:  false,
		Destroyer:  false,
	}, hitShipAt: make(map[[2]int]uint8)}
}

// **FIXED**: This function is now a safe, read-only check of the board's state.
// It no longer modifies the board's data, preventing the bug.
func (b *battleshipBoard) IsShipSunk(ship uint8) (sunk bool, err error) {
	sunk = b.sunkShips[ship]
	return sunk, nil
}

// **NEW**: This internal method contains the logic to scan the board and update the sunk status.
// It is only called from within the Attack function.
func (b *battleshipBoard) updateSunkStatus(ship uint8) (sunk bool, err error) {
	// If already marked as sunk, no need to check again.
	if b.sunkShips[ship] {
		return true, nil
	}
	// Scan the board for any remaining pieces of the ship.
	for _, item := range b.FlatSlice() {
		if item == ship {
			return false, nil // Found a piece, so it's not sunk.
		}
	}
	// No pieces were found, so the ship is now sunk.
	b.sunkShips[ship] = true
	return true, nil
}

func (b *battleshipBoard) CopyHitValues(otherBoard BattleshipBoard) {
	for rows := 0; rows < otherBoard.Rows(); rows++ {
		for cols := 0; cols < otherBoard.Cols(); cols++ {
			if b.Coordinate(rows, cols) <= 3 {
				b.SetCoordinate(rows, cols, otherBoard.Coordinate(rows, cols))
			}
		}
	}
}

func (b *battleshipBoard) Attack(x, y int) (hit, sunk bool, shipType uint8, err error) {
	currentValue := b.Coordinate(x, y)
	if currentValue == Hit || currentValue == Miss {
		return false, false, 0, errors.New("already hit or missed at this location")
	}
	if currentValue == Empty {
		b.SetCoordinate(x, y, Miss)
		return false, false, 0, nil
	}

	b.hitShipAt[[2]int{x, y}] = currentValue
	b.SetCoordinate(x, y, Hit)

	// Check if this hit caused the ship to sink.
	sunk, err = b.updateSunkStatus(currentValue)
	if err != nil {
		return false, false, 0, err
	}

	// If the ship is sunk, update all its hit cells to the SUNK state.
	if sunk {
		for coord, ship := range b.hitShipAt {
			if ship == currentValue {
				b.SetCoordinate(coord[0], coord[1], SUNK)
			}
		}
	}

	return true, sunk, currentValue, nil
}

func (b *battleshipBoard) RecordSunkShip(shipType uint8) {
	if _, ok := b.sunkShips[shipType]; !ok {
		b.sunkShips[shipType] = true
	}
}

// ... (The rest of the file remains the same) ...
func (b *battleshipBoard) PlaceShip(x int, y int, shipType uint8, orientation uint8) bool {
	// Determine the length of the ship from the ship type
	length := ShipLength(shipType)
	if canPlace := b.CanPlace(x, y, length, orientation); !canPlace {
		//fmt.Println("Warning: cannot place ship at", x, y, "with length", length, "and orientation", orientation)
		return false
	}
	// Normalize orientation for simple handling
	switch orientation {

	case Horizontal:
		// Place to the right from (x, y)
		for i := 0; i < length; i++ {
			// bounds check to avoid panics on overflow
			if x+i < b.Cols() && y < b.Rows() {
				b.SetCoordinate(x+i, y, shipType)
			}
		}

	case Vertical:
		// Place downward from (x, y)
		for i := 0; i < length; i++ {
			// bounds check to avoid panics on overflow
			if x < b.Cols() && y+i < b.Rows() {
				b.SetCoordinate(x, y+i, shipType)
			}
		}
	default:
		fmt.Println("Unknown orientation: ", orientation)
		return false
		// Unknown orientation; do nothing or optionally log
	}

	return true
}

func (b *battleshipBoard) SeedBoard() {
	// reset the board to empty
	for x := 0; x < b.Cols(); x++ {
		for y := 0; y < b.Rows(); y++ {
			b.SetCoordinate(x, y, 0)
		}
	}
	// reset sunkShips and hitShipAt tracking
	b.sunkShips = map[uint8]bool{
		Carrier: false, Battleship: false, Cruiser: false, Submarine: false, Destroyer: false,
	}
	if b.hitShipAt == nil {
		b.hitShipAt = make(map[[2]int]uint8)
	} else {
		for k := range b.hitShipAt {
			delete(b.hitShipAt, k)
		}
	}
	// place ships randomly
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shipTypes := []uint8{Carrier, Battleship, Cruiser, Submarine, Destroyer}
	for _, shipType := range shipTypes {
		placed := false
		// try up to a reasonable number of attempts
		for attempts := 0; attempts < 1000 && !placed; attempts++ {
			orientation := uint8(r.Intn(2))
			var x, y int
			x = rand.Intn(b.Cols())
			y = rand.Intn(b.Rows())
			placed = b.PlaceShip(x, y, shipType, orientation)
		}
		if !placed {
			fmt.Println("Warning: could not place ship type", shipType)
		}
	}
}

func (b *battleshipBoard) CanPlace(x, y, length int, orientation uint8) bool {
	if orientation == Horizontal {
		if x+length > b.Cols() {
			return false
		}
		for i := 0; i < length; i++ {
			if b.Coordinate(x+i, y) != Empty {
				return false
			}
		}
	} else if orientation == Vertical {
		if y+length > b.Rows() {
			return false
		}
		for i := 0; i < length; i++ {
			if b.Coordinate(x, y+i) != Empty {
				return false
			}
		}
	}
	return true
}

func ShipLength(shipType uint8) int {
	switch shipType {
	case Carrier:
		return 5
	case Battleship:
		return 4
	case Cruiser:
		return 3
	case Submarine:
		return 3
	case Destroyer:
		return 2
	default:
		return 0
	}
}

func (b *battleshipBoard) IsCellSunk(x, y int) bool {
	coordState := b.Coordinate(x, y)
	if coordState != Hit && coordState != SUNK {
		return false
	}
	ship, ok := b.hitShipAt[[2]int{x, y}]
	if !ok {
		return coordState == SUNK
	}
	sunk, _ := b.IsShipSunk(ship)
	return sunk
}

func (b *battleshipBoard) AllShipsSunk() bool {
	shipTypes := []uint8{Carrier, Battleship, Cruiser, Submarine, Destroyer}
	for _, st := range shipTypes {
		sunk, _ := b.IsShipSunk(st)
		if !sunk {
			return false
		}
	}
	return true
}

func (b *battleshipBoard) SunkShips() map[uint8]bool {
	return b.sunkShips
}

func (b *battleshipBoard) HitShipAt() map[[2]int]uint8 {
	return b.hitShipAt
}

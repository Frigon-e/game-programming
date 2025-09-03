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
	Carrier    uint8 = 3
	Battleship uint8 = 4
	Cruiser    uint8 = 5
	Submarine  uint8 = 6
	Destroyer  uint8 = 7
)

const (
	Horizontal uint8 = 0
	Vertical   uint8 = 1
)

type BattleshipBoard interface {
	ddd.Board[uint8]
	SeedBoard()
	Attack(x, y int) (hit, sunk bool, err error)
}

type battleshipBoard struct {
	ddd.Board[uint8] // embed the generic board to reuse its methods
	sunkShips        map[uint8]bool
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
	}}
}

func (b *battleshipBoard) PlaceShip(x int, y int, shipType uint8, orientation uint8) bool {
	// Determine the length of the ship from the ship type
	length := shipLength(shipType)
	if canPlace := b.CanPlace(x, y, length, orientation); !canPlace {
		fmt.Println("Warning: cannot place ship at", x, y, "with length", length, "and orientation", orientation)
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
	// place ships randomly
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shipTypes := []uint8{Carrier, Battleship, Cruiser, Submarine, Destroyer}
	for _, shipType := range shipTypes {
		placed := false
		// try up to a reasonable number of attempts
		for attempts := 0; attempts < 1000 && !placed; attempts++ {
			orientation := uint8(r.Intn(2)) // 0 horizontal, 1 vertical
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
	if orientation == Horizontal { // horizontal
		if x+length >= b.Cols() {
			return false
		}
		for i := 0; i < length; i++ {
			if b.Coordinate(x+i, y) != Empty {
				return false
			}
		}
	} else if orientation == Vertical {
		if y+length >= b.Rows() {
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

func (b *battleshipBoard) IsShipSunk(ship uint8) (sunk bool, err error) {
	if b.sunkShips[ship] {
		return true, nil
	}
	for _, item := range b.FlatSlice() {
		if item == ship {
			return false, nil
		}
	}
	b.sunkShips[ship] = true
	return true, nil

}

func (b *battleshipBoard) Attack(x, y int) (hit, sunk bool, err error) {
	currentValue := b.Coordinate(x, y)
	if currentValue == Hit || currentValue == Miss {
		return false, false, errors.New("already hit or missed at this location")
	}
	if currentValue == Empty {
		b.SetCoordinate(x, y, Miss)
		return false, false, nil
	}
	b.SetCoordinate(x, y, Hit)
	sunk, err = b.IsShipSunk(currentValue)
	if err != nil {
		return false, false, err
	}
	return true, sunk, nil
}

func shipLength(shipType uint8) int {
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

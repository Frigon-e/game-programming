package ddd

import (
	"SideProjectGames/internal/ddd"
	"math/rand"
	"time"
)

// GolBoard composes the generic ddd.Board with extra Game of Life helpers.
// In Go, we "inherit" behavior by embedding interfaces/structs (composition),
// not by classical inheritance.
type GolBoard interface {
	ddd.Board[bool]
	SeedBoard()
	CountSurroundingLive(x int, y int) int
}

type golBoard struct {
	ddd.Board[bool] // embed the generic board to reuse its methods
}

var _ GolBoard = (*golBoard)(nil)

func NewGOLBoard(width int, height int) GolBoard {
	return newGOLBoard(width, height)
}

func newGOLBoard(width int, height int) GolBoard {
	return &golBoard{Board: ddd.NewBoard[bool](width, height)}
}

func (b *golBoard) CountSurroundingLive(x int, y int) int {
	surroundArray := []int{-1, 0, 1}
	totalAlive := 0
	for _, rowOffset := range surroundArray {
		for _, colOffset := range surroundArray {
			if rowOffset == 0 && colOffset == 0 {
				continue
			}

			if b.Coordinate(x+rowOffset, y+colOffset) {
				totalAlive += 1
			}
		}
	}

	return totalAlive
}

func (b *golBoard) SeedBoard() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Work on a copy then write back via CopyBoard to respect encapsulation
	slice := b.FlatSlice()
	for i := 0; i < len(slice); i++ {
		// Rough 50/50 seed; adjust as desired
		slice[i] = r.Intn(2) == 1
	}
	b.CopyBoard(slice)
}

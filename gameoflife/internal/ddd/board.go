package ddd

import (
	"SideProjectGames/internal/ddd"
	"math/rand"
	"time"
)

// GolBoard composes the generic ddd.Board with extra Game of Life helpers.
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
	for _, yOffset := range surroundArray {
		for _, xOffset := range surroundArray {
			if yOffset == 0 && xOffset == 0 {
				continue
			}

			if b.Coordinate(x+xOffset, y+yOffset) {
				totalAlive += 1
			}
		}
	}

	return totalAlive
}

func (b *golBoard) SeedBoard() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	slice := b.FlatSlice()
	for i := 0; i < len(slice); i++ {
		slice[i] = r.Intn(2) == 1
	}
	b.CopyBoard(slice)
}

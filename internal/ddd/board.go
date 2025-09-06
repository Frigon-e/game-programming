package ddd

import "fmt"

type Board[T any] interface {
	Coordinate(x int, y int) T
	SetCoordinate(x int, y int, value T)
	CopyBoard(flatSlice []T)
	PrintBoard()
	FlatSlice() []T
	Rows() int
	Cols() int
}

type board[T any] struct {
	flatSlice []T
	cols      int
	rows      int
	Board[T]
}

var _ Board[int] = (*board[int])(nil)

func NewBoard[T any](width int, height int) Board[T] {
	return newBoard[T](width, height)
}

func newBoard[T any](width int, height int) Board[T] {
	return &board[T]{
		rows:      height,
		cols:      width,
		flatSlice: make([]T, width*height),
	}
}

func (b *board[T]) CopyBoard(writeSlice []T) {
	copy(b.flatSlice, writeSlice)
}

func (b *board[T]) Coordinate(x int, y int) T {
	row := (y + b.rows) % b.rows
	col := (x + b.cols) % b.cols

	return b.flatSlice[(row*b.cols)+col]
}

func (b *board[T]) SetCoordinate(x int, y int, value T) {
	row := (y + b.rows) % b.rows
	col := (x + b.cols) % b.cols

	b.flatSlice[(row*b.cols)+col] = value
}

func (b *board[T]) PrintBoard() {
	for y := 0; y < b.rows; y++ {
		for x := 0; x < b.cols; x++ {
			fmt.Print(" ", b.Coordinate(x, y), " ")
		}
		fmt.Println()
	}
}

func (b *board[T]) FlatSlice() []T {
	return b.flatSlice
}

func (b *board[T]) Cols() int {
	return b.cols
}

func (b *board[T]) Rows() int {
	return b.rows
}

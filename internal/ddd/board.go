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
	for index := 0; index < len(b.flatSlice); index++ {
		b.flatSlice[index] = writeSlice[index]
	}
}

func (b *board[T]) Coordinate(x int, y int) T {
	newRow := (x + b.rows) % b.rows
	newCol := (y + b.cols) % b.cols

	return b.flatSlice[(newRow*b.rows)+newCol]
}

func (b *board[T]) SetCoordinate(x int, y int, value T) {
	newRow := (x + b.rows) % b.rows
	newCol := (y + b.cols) % b.cols

	b.flatSlice[(newRow*b.rows)+newCol] = value
}

func (b *board[T]) PrintBoard() {
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			fmt.Print(b.Coordinate(row, col))
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

package application

import "math/rand"

type AI interface {
	TakeTurn(board BattleshipBoard) (x, y int)
}

type simpleAI struct {
	rows int
	cols int
}

func NewSimpleAI(rows, cols int) AI {
	return &simpleAI{
		rows: rows,
		cols: cols,
	}
}

func (ai *simpleAI) TakeTurn(board BattleshipBoard) (x, y int) {
	// For now, the AI will just make a random move.
	// A more advanced AI would keep track of previous moves.
	x = rand.Intn(ai.cols)
	y = rand.Intn(ai.rows)
	return x, y
}

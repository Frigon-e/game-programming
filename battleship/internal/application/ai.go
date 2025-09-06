package application

import (
	"math/rand"
)

// TakeTurn implements the AI interface, generating a move based on a heatmap.
func TakeTurn(board BattleshipBoard) (x, y int) {
	// 1. Create a new heatmap for this turn.
	heatMap := NewHeatmapBoard(board.Cols(), board.Rows())

	// 2. Calculate the heatmap based on the current state of the opponent's board.
	heatMap.CalculateHeatmap(board)

	//heatMap.PrintBoard()

	// 3. Find the coordinate(s) with the highest "heat" that haven't been attacked yet.
	var bestCoords [][2]int
	var maxHeat int16 = -1

	for y := 0; y < board.Rows(); y++ {
		for x := 0; x < board.Cols(); x++ {
			// We only consider empty cells as potential targets.
			if board.Coordinate(x, y) == Empty {
				currentHeat := heatMap.Coordinate(x, y)
				if currentHeat > maxHeat {
					maxHeat = currentHeat
					bestCoords = [][2]int{{x, y}} // Start a new list of best coordinates
				} else if currentHeat == maxHeat {
					bestCoords = append(bestCoords, [2]int{x, y}) // Add to the list of best coordinates
				}
			}
		}
	}

	// 4. If high-value targets are found, pick one at random from the best options.
	if len(bestCoords) > 0 {
		choice := rand.Intn(len(bestCoords))
		return bestCoords[choice][0], bestCoords[choice][1]
	}

	// 5. If no high-value targets are found (e.g., on the first turn),
	// pick a random valid spot as a fallback.
	for {
		randX := rand.Intn(board.Cols())
		randY := rand.Intn(board.Rows())
		if board.Coordinate(randX, randY) == Empty {
			return randX, randY
		}
	}
}

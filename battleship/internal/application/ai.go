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

	// 3. Find the coordinate(s) with the highest "heat" that hasn't been attacked yet.
	bestCoords := heatMap.GetBestCoords(board)

	// 4. If there is a tiebreaker make sure to implement the code
	newBestCoords := make([][2]int, 0, len(bestCoords))
	var maxHeat int16 = -1
	if len(bestCoords) >= 2 {
		for _, coord := range bestCoords {
			currentHeat := heatMap.SumNeighbours(coord[0], coord[1])
			if currentHeat > maxHeat {
				maxHeat = currentHeat
				newBestCoords = [][2]int{{coord[0], coord[1]}} // Start a new list of best coordinates
			} else if currentHeat == maxHeat {
				newBestCoords = append(newBestCoords, [2]int{coord[0], coord[1]}) // Add to the list of best coordinates
			}
		}
	} else if len(bestCoords) == 1 {
		return bestCoords[0][0], bestCoords[0][1]
	}

	// 5. If high-value targets are found, use it
	if len(newBestCoords) > 0 {
		choice := rand.Intn(len(newBestCoords))
		return newBestCoords[choice][0], newBestCoords[choice][1]
	}

	// 6. If no high-value targets are found (e.g., on the first turn),
	// pick a random valid spot as a fallback.
	for {
		randX := rand.Intn(board.Cols())
		randY := rand.Intn(board.Rows())
		if board.Coordinate(randX, randY) == Empty {
			return randX, randY
		}
	}
}

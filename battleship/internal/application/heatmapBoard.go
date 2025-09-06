package application

import "SideProjectGames/internal/ddd"

type HeatmapBoard interface {
	ddd.Board[int16]
	CalculateHeatmap(bsboard BattleshipBoard)
}

type heatmapBoard struct {
	ddd.Board[int16]
}

var _ HeatmapBoard = (*heatmapBoard)(nil)

func NewHeatmapBoard(width, height int) HeatmapBoard {
	return newHeatmapBoard(width, height)
}

func newHeatmapBoard(width, height int) HeatmapBoard {
	return &heatmapBoard{Board: ddd.NewBoard[int16](width, height)}
}

// CalculateHeatmap generates the probability map for the AI to make a decision.
func (hm *heatmapBoard) CalculateHeatmap(bsBoard BattleshipBoard) {
	// 1. Reset the heatmap to all zeros before recalculation.
	for i := 0; i < hm.Rows()*hm.Cols(); i++ {
		hm.FlatSlice()[i] = 0
	}

	// 2. Determine which ships are still alive.
	allShipTypes := []uint8{Carrier, Battleship, Cruiser, Submarine, Destroyer}
	aliveShipLengths := []int{}
	for _, shipType := range allShipTypes {
		sunk, _ := bsBoard.IsShipSunk(shipType)
		if !sunk {
			aliveShipLengths = append(aliveShipLengths, ShipLength(shipType))
		}
	}

	// 3. TARGET MODE: Calculate base probabilities.
	// Iterate over every cell and every alive ship to see how many valid placements exist.
	for _, length := range aliveShipLengths {
		for r := 0; r < hm.Rows(); r++ {
			for c := 0; c < hm.Cols(); c++ {
				// Horizontal check: Can a ship of this length be placed here horizontally?
				if canPlaceShip(bsBoard, c, r, length, Horizontal) {
					// If yes, increment the heat for all cells this ship would occupy.
					for i := 0; i < length; i++ {
						hm.SetCoordinate(c+i, r, hm.Coordinate(c+i, r)+1)
					}
				}
				// Vertical check: Can a ship of this length be placed here vertically?
				if canPlaceShip(bsBoard, c, r, length, Vertical) {
					// If yes, increment the heat.
					for i := 0; i < length; i++ {
						hm.SetCoordinate(c, r+i, hm.Coordinate(c, r+i)+1)
					}
				}
			}
		}
	}

	// 4. HUNT MODE: Add high-value heat around existing hits.
	// This prioritizes sinking an already-discovered ship over finding a new one.
	for r := 0; r < hm.Rows(); r++ {
		for c := 0; c < hm.Cols(); c++ {
			if bsBoard.Coordinate(c, r) == Hit {
				// Add a large bonus to adjacent squares to encourage "hunting."
				// Ensure you don't add heat to spots that are already missed.
				if r-1 >= 0 && bsBoard.Coordinate(c, r-1) == Empty {
					hm.SetCoordinate(c, r-1, hm.Coordinate(c, r-1)+100)
				}
				if r+1 < hm.Rows() && bsBoard.Coordinate(c, r+1) == Empty {
					hm.SetCoordinate(c, r+1, hm.Coordinate(c, r+1)+100)
				}
				if c-1 >= 0 && bsBoard.Coordinate(c-1, r) == Empty {
					hm.SetCoordinate(c-1, r, hm.Coordinate(c-1, r)+100)
				}
				if c+1 < hm.Cols() && bsBoard.Coordinate(c+1, r) == Empty {
					hm.SetCoordinate(c+1, r, hm.Coordinate(c+1, r)+100)
				}
			}
		}
	}
}

// canPlaceShip is a helper to check if a ship can be placed without overlapping misses.
// This is used for heatmap generation, not for initial board seeding.
func canPlaceShip(board BattleshipBoard, x, y, length int, orientation uint8) bool {
	if orientation == Horizontal {
		if x+length > board.Cols() {
			return false
		}
		for i := 0; i < length; i++ {
			// A placement is invalid if it overlaps a Miss.
			if board.Coordinate(x+i, y) == Miss {
				return false
			}
		}
	} else { // Vertical
		if y+length > board.Rows() {
			return false
		}
		for i := 0; i < length; i++ {
			if board.Coordinate(x, y+i) == Miss {
				return false
			}
		}
	}
	return true
}

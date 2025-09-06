package application

import (
	"SideProjectGames/internal/ddd"
)

type HeatmapBoard interface {
	ddd.Board[int16]
	CalculateHeatmap(bsboard BattleshipBoard)
	SumNeighbours(x int, y int) int16
	GetBestCoords(board BattleshipBoard) [][2]int
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
	aliveShipLengths := []uint8{}
	// Prefer using the sunk ships map directly to avoid any side effects or stale state from IsShipSunk.
	for shipLength, isShipSunk := range bsBoard.SunkShips() {
		if !isShipSunk {
			aliveShipLengths = append(aliveShipLengths, uint8(ShipLength(shipLength)))
		}
	}

	// 3. TARGET MODE: Calculate base probabilities.
	// Iterate over every cell and every alive ship to see how many valid placements exist.
	for _, length := range aliveShipLengths {
		intShipLength := int(length)
		for r := 0; r < hm.Rows(); r++ {
			for c := 0; c < hm.Cols(); c++ {

				// Horizontal check: Can a ship of this length be placed here horizontally?
				if canPlaceShip(bsBoard, c, r, intShipLength, Horizontal) {
					// If yes, increment the heat for all cells this ship would occupy.
					for i := 0; i < intShipLength; i++ {
						hm.SetCoordinate(c+i, r, hm.Coordinate(c+i, r)+1)
					}
				}
				// Vertical check: Can a ship of this length be placed here vertically?
				if canPlaceShip(bsBoard, c, r, intShipLength, Vertical) {
					// If yes, increment the heat.
					for i := 0; i < intShipLength; i++ {
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
			if bsBoard.Coordinate(c, r) == Empty {
				isVerticalGap := (r > 0 && bsBoard.Coordinate(c, r-1) == Hit) && (r < hm.Rows()-1 && bsBoard.Coordinate(c, r+1) == Hit)
				isHorizontalGap := (c > 0 && bsBoard.Coordinate(c-1, r) == Hit) && (c < hm.Cols()-1 && bsBoard.Coordinate(c+1, r) == Hit)
				bonus := int16(500)

				if isVerticalGap || isHorizontalGap {
					hm.SetCoordinate(c, r, hm.Coordinate(c, r)+bonus)
				}
			}

			if bsBoard.Coordinate(c, r) == Hit {
				// Check for adjacent hits to determine if we've found a line.
				// Adds Extra Bonus in case of multiple Hits in a row
				verticalBonus := int16(0)
				if r > 0 && bsBoard.Coordinate(c, r-1) == Hit {
					verticalBonus += 500
				}
				if r < hm.Rows()-1 && bsBoard.Coordinate(c, r+1) == Hit {
					verticalBonus += 500
				}

				// Adds Extra Bonus in case of multiple Hits in a row
				horizontalBonus := int16(0)
				if c > 0 && bsBoard.Coordinate(c-1, r) == Hit {
					horizontalBonus += 500
				}
				if c < hm.Cols()-1 && bsBoard.Coordinate(c+1, r) == Hit {
					horizontalBonus += 500
				}

				// North neighbour
				if r > 0 && bsBoard.Coordinate(c, r-1) == Empty {
					bonus := int16(100) + verticalBonus
					hm.SetCoordinate(c, r-1, hm.Coordinate(c, r-1)+bonus)
				}
				// South neighbour
				if r < hm.Rows()-1 && bsBoard.Coordinate(c, r+1) == Empty {
					bonus := int16(100) + verticalBonus
					hm.SetCoordinate(c, r+1, hm.Coordinate(c, r+1)+bonus)
				}
				// West neighbour
				if c > 0 && bsBoard.Coordinate(c-1, r) == Empty {
					bonus := int16(100) + horizontalBonus
					if verticalBonus > 0 && horizontalBonus == 0 {
						bonus = 0
					}
					hm.SetCoordinate(c-1, r, hm.Coordinate(c-1, r)+bonus)
				}
				// East neighbour
				if c < hm.Cols()-1 && bsBoard.Coordinate(c+1, r) == Empty {
					bonus := int16(100) + horizontalBonus
					if verticalBonus > 0 && horizontalBonus == 0 {
						bonus = 0
					}
					hm.SetCoordinate(c+1, r, hm.Coordinate(c+1, r)+bonus)
				}
			}
		}
	}
}

func (hm *heatmapBoard) SumNeighbours(x, y int) int16 {
	surroundArray := []int{-1, 0, 1}
	totalSum := int16(0)
	for _, yOffset := range surroundArray {
		for _, xOffset := range surroundArray {
			if yOffset == 0 && xOffset == 0 {
				continue
			}

			if x+xOffset < 0 || x+xOffset > hm.Cols()-1 {
				continue
			}

			if y+yOffset < 0 || y+yOffset > hm.Rows()-1 {
				continue
			}

			totalSum += hm.Coordinate(x+xOffset, y+yOffset)
		}
	}

	return totalSum
}

func (hm *heatmapBoard) GetBestCoords(board BattleshipBoard) [][2]int {
	var bestCoords [][2]int
	var maxHeat int16 = -1

	for y := 0; y < board.Rows(); y++ {
		for x := 0; x < board.Cols(); x++ {
			// We only consider empty cells as potential targets.
			if board.Coordinate(x, y) == Empty {
				currentHeat := hm.Coordinate(x, y)
				if currentHeat > maxHeat {
					maxHeat = currentHeat
					bestCoords = [][2]int{{x, y}} // Start a new list of best coordinates
				} else if currentHeat == maxHeat {
					bestCoords = append(bestCoords, [2]int{x, y}) // Add to the list of best coordinates
				}
			}
		}
	}

	return bestCoords
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
			if board.Coordinate(x+i, y) != Empty {
				return false
			}
		}
	} else { // Vertical
		if y+length > board.Rows() {
			return false
		}
		for i := 0; i < length; i++ {
			if board.Coordinate(x, y+i) != Empty {
				return false
			}
		}
	}
	return true
}

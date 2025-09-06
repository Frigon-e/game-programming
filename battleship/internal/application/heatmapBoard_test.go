package application

import "testing"

func TestCalculateHeatmap_PopulatesOnEmptyViewBoard(t *testing.T) {
	// Given an empty AI view board (no ships placed, no shots taken)
	view := NewBattleshipBoard(10, 10)
	// Do not call SeedBoard on the view board; it should start empty with all ships "unknown/not sunk".

	hm := NewHeatmapBoard(10, 10)
	hm.CalculateHeatmap(view)

	// Then the heatmap should have non-zero values since many placements are possible.
	var total int64
	for _, v := range hm.FlatSlice() {
		total += int64(v)
	}
	if total == 0 {
		t.Fatalf("expected non-zero heat on an empty board, got total=%d", total)
	}
}

func TestCalculateHeatmap_HeatReducesWhenAShipIsSunk(t *testing.T) {
	view := NewBattleshipBoard(10, 10)

	hm := NewHeatmapBoard(10, 10)
	hm.CalculateHeatmap(view)
	var baseline int64
	for _, v := range hm.FlatSlice() {
		baseline += int64(v)
	}
	if baseline == 0 {
		t.Fatalf("expected baseline heat > 0, got %d", baseline)
	}

	// Mark the smallest ship (Destroyer) as sunk and recalc
	view.RecordSunkShip(Destroyer)
	hm.CalculateHeatmap(view)
	var after int64
	for _, v := range hm.FlatSlice() {
		after += int64(v)
	}
	if !(after < baseline) {
		t.Fatalf("expected total heat to decrease after sinking Destroyer; baseline=%d, after=%d", baseline, after)
	}
}

func TestCalculateHeatmap_ZeroWhenAllShipsSunk_NoHits(t *testing.T) {
	view := NewBattleshipBoard(10, 10)
	view.RecordSunkShip(Carrier)
	view.RecordSunkShip(Battleship)
	view.RecordSunkShip(Cruiser)
	view.RecordSunkShip(Submarine)
	view.RecordSunkShip(Destroyer)

	hm := NewHeatmapBoard(10, 10)
	hm.CalculateHeatmap(view)
	var total int64
	for _, v := range hm.FlatSlice() {
		total += int64(v)
	}
	if total != 0 {
		t.Fatalf("expected zero heat when all ships are sunk and no hits, got total=%d", total)
	}
}

func TestCalculateHeatmap_HuntModeAddsHeatEvenIfAllSunk(t *testing.T) {
	view := NewBattleshipBoard(10, 10)
	// Mark all ships as sunk
	view.RecordSunkShip(Carrier)
	view.RecordSunkShip(Battleship)
	view.RecordSunkShip(Cruiser)
	view.RecordSunkShip(Submarine)
	view.RecordSunkShip(Destroyer)
	// Place a hit on the board to trigger hunt mode bonuses
	view.SetCoordinate(5, 5, Hit)

	hm := NewHeatmapBoard(10, 10)
	hm.CalculateHeatmap(view)
	var total int64
	for _, v := range hm.FlatSlice() {
		total += int64(v)
	}
	if total == 0 {
		t.Fatalf("expected non-zero heat due to hunt mode bonuses around a hit, got total=%d", total)
	}
}

func TestTakeTurn_TiebreakerBySumSurrounding(t *testing.T) {
	// 1. Arrange
	// Create a board state that will result in a tie between multiple coordinates.
	// We want to ensure that the tie-breaking logic, which uses the sum of
	// surrounding heat, correctly picks the best spot.
	board := NewBattleshipBoard(10, 10)

	// Place two separate hits to create two hotspots with equally valuable targets.
	board.SetCoordinate(1, 1, Hit)
	board.SetCoordinate(8, 8, Hit)

	// Now, create an asymmetry to test the tie-breaker. We'll place a 'Miss'
	// near one hotspot. This reduces the sum of neighboring heat for that spot,
	// making it less attractive after the tie-breaker logic is applied.
	board.SetCoordinate(0, 3, Miss)
	board.SetCoordinate(3, 0, Miss)

	// 2. Act
	x, y := TakeTurn(board)

	// 3. Assert
	// This is where it should hit if it's working properly
	expectedWinnerCoords := [][2]int{
		{7, 8}, {9, 8}, {8, 7}, {8, 9},
	}

	isWinner := false
	for _, coord := range expectedWinnerCoords {
		if x == coord[0] && y == coord[1] {
			isWinner = true
			break
		}
	}

	if !isWinner {
		t.Errorf("TakeTurn() chose (%d, %d), which was not in the expected set of winning coordinates %v", x, y, expectedWinnerCoords)
	}

	// These show that it doesn't work and it's not adding properly
	losingCoords := [][2]int{
		{1, 0}, {0, 1}, {2, 1}, {1, 2},
	}

	for _, coord := range losingCoords {
		if x == coord[0] && y == coord[1] {
			t.Errorf("TakeTurn() chose (%d, %d), which should have been eliminated by the tie-breaker", x, y)
		}
	}
}

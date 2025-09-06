package battleship

import (
	"SideProjectGames/battleship/internal/application"
	"SideProjectGames/internal/config"
	"bytes"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

// Run launches an Ebiten window to visualize a Battleship game loop using the provided config.
// This mirrors the Game of Life loop structure and prepares for separate User and AI boards.
func Run(cfg config.AppConfig) error {
	g := &game{
		cellSize:        50,
		stepEvery:       time.Millisecond * 100, // kept for consistency; not used yet for turn timing
		rows:            cfg.BATTLESHIPHEIGHT,
		cols:            cfg.BATTLESHIPWIDTH,
		aiSolutionBoard: application.NewBattleshipBoard(cfg.BATTLESHIPWIDTH, cfg.BATTLESHIPHEIGHT),
		userBoard:       application.NewBattleshipBoard(cfg.BATTLESHIPWIDTH, cfg.BATTLESHIPHEIGHT),
		aiViewBoard:     application.NewBattleshipBoard(cfg.BATTLESHIPWIDTH, cfg.BATTLESHIPHEIGHT),
		isPlayerTurn:    true,
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	g.userBoard.SeedBoard()
	g.userBoard.PrintBoard()

	g.aiSolutionBoard.SeedBoard()
	// Layout: two boards stacked vertically with a gap
	gap := 20
	boardW := g.cols * g.cellSize
	boardH := g.rows * g.cellSize
	w := boardW
	h := boardH*2 + gap
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Battleship")

	return ebiten.RunGame(g)
}

type game struct {
	rows, cols      int
	cellSize        int
	stepEvery       time.Duration
	lastStep        time.Time
	userBoard       application.BattleshipBoard
	aiSolutionBoard application.BattleshipBoard
	aiViewBoard     application.BattleshipBoard
	isPlayerTurn    bool
	gameOver        bool
	winner          string
}

func (g *game) Update() error {
	if g.gameOver {
		return nil
	}
	if g.isPlayerTurn {
		g.handleClick()
	} else {
		g.step()
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// Colors
	bgColor := color.RGBA{R: 20, G: 30, B: 40, A: 255}
	gridColor := color.RGBA{R: 160, G: 170, B: 180, A: 255}
	hitColor := color.RGBA{R: 255, G: 100, B: 30, A: 255}
	missColor := color.RGBA{R: 40, G: 60, B: 80, A: 255}
	sunkColor := color.RGBA{R: 180, G: 30, B: 180, A: 255}

	// Clear
	screen.Fill(bgColor)

	cs := g.cellSize
	boardH := g.rows * cs
	gap := 20

	// Determine opacity based on whose turn it is
	userBoardAlpha := uint8(255)
	aiBoardAlpha := uint8(255)
	if g.isPlayerTurn {
		aiBoardAlpha = 128 // Dim the AI board if it's the player's turn
	} else {
		userBoardAlpha = 128 // Dim the user board if it's the AI's turn
	}

	// Draw user board (top)
	drawGrid(screen, 0, 0, g.cols, g.rows, cs, applyAlpha(gridColor, userBoardAlpha))
	for x := 0; x < g.cols; x++ {
		for y := 0; y < g.rows; y++ {
			var chosenColor color.Color
			cell := g.aiSolutionBoard.Coordinate(x, y)
			switch cell {
			case application.Hit:
				if g.aiSolutionBoard.IsCellSunk(x, y) {
					chosenColor = sunkColor
				} else {
					chosenColor = hitColor
				}
			case application.Miss:
				chosenColor = missColor
			default:
				chosenColor = bgColor
			}
			xPix := x*cs + 1
			yPix := y*cs + 1
			vector.DrawFilledRect(screen, float32(xPix), float32(yPix), float32(cs-2), float32(cs-2), applyAlpha(chosenColor, userBoardAlpha), false)
			// Draw X overlay for sunk cells
			if cell == application.SUNK {
				lineCol := applyAlpha(color.RGBA{R: 255, G: 255, B: 255, A: 255}, userBoardAlpha)
				x1 := float32(xPix + 2)
				y1 := float32(yPix + 2)
				x2 := float32(xPix + cs - 3)
				y2 := float32(yPix + cs - 3)
				vector.StrokeLine(screen, x1, y1, x2, y2, 2, lineCol, false)
				vector.StrokeLine(screen, x1, y2, x2, y1, 2, lineCol, false)
			}
		}
	}

	// Draw AI board (bottom)
	offsetY := boardH + gap
	drawGrid(screen, 0, offsetY, g.cols, g.rows, cs, applyAlpha(gridColor, aiBoardAlpha))
	for x := 0; x < g.cols; x++ {
		for y := 0; y < g.rows; y++ {
			var chosenColor color.Color
			cell := g.userBoard.Coordinate(x, y)
			switch cell {
			case application.Hit, application.SUNK:
				if g.userBoard.IsCellSunk(x, y) {
					chosenColor = sunkColor
				} else {
					chosenColor = hitColor
				}
			case application.Miss:
				chosenColor = missColor
			default:
				chosenColor = bgColor
			}
			xPix := x*cs + 1
			yPix := offsetY + y*cs + 1
			vector.DrawFilledRect(screen, float32(xPix), float32(yPix), float32(cs-2), float32(cs-2), applyAlpha(chosenColor, aiBoardAlpha), false)
			// Draw X overlay for sunk cells
			if cell == application.SUNK {
				lineCol := applyAlpha(color.RGBA{R: 255, G: 255, B: 255, A: 255}, aiBoardAlpha)
				x1 := float32(xPix + 2)
				y1 := float32(yPix + 2)
				x2 := float32(xPix + cs - 3)
				y2 := float32(yPix + cs - 3)
				vector.StrokeLine(screen, x1, y1, x2, y2, 2, lineCol, false)
				vector.StrokeLine(screen, x1, y2, x2, y1, 2, lineCol, false)
			}
		}
	}

	// UI text
	msg := fmt.Sprintf("User board (top) | AI board (bottom)   Cells: %dx%d  CellSize: %d", g.cols, g.rows, g.cellSize)
	op := &text.DrawOptions{}
	op.GeoM.Translate(10, float64(g.rows*g.cellSize*2+gap-28))
	op.ColorScale.ScaleWithColor(color.RGBA{0, 0, 0, 255})
	text.Draw(screen, msg, &text.GoTextFace{Source: mplusFaceSource, Size: 18}, op)

	// Game over message
	if g.gameOver {
		winnerMsg := fmt.Sprintf("Game Over - %s wins!", g.winner)
		op2 := &text.DrawOptions{}
		op2.GeoM.Translate(10, float64(g.rows*g.cellSize*2+gap-56))
		op2.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
		text.Draw(screen, winnerMsg, &text.GoTextFace{Source: mplusFaceSource, Size: 24}, op2)
	}
}

func drawGrid(screen *ebiten.Image, offsetX, offsetY, cols, rows, cellSize int, col color.Color) {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			xPix := offsetX + x*cellSize
			yPix := offsetY + y*cellSize
			vector.DrawFilledRect(screen, float32(xPix), float32(yPix), float32(cellSize-1), float32(cellSize-1), col, false)
		}
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *game) handleClick() {
	if g.gameOver || !g.isPlayerTurn {
		return // Ignore clicks if it's not the player's turn or game is over
	}
	mouseX, mouseY := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cs := g.cellSize
		boardW := g.cols * cs
		boardH := g.rows * cs
		gap := 20
		offsetY := boardH + gap

		// Check if click is on AI board (bottom)
		if mouseX >= 0 && mouseX < boardW && mouseY >= offsetY && mouseY < offsetY+boardH {
			gridX := mouseX / cs
			gridY := (mouseY - offsetY) / cs
			// Toggle a placeholder shot state; hit simulation toggles true/false
			hit, sunk, _, err := g.userBoard.Attack(gridX, gridY)
			if err != nil {
				return // Already clicked here
			}
			if !hit {
				g.isPlayerTurn = false
				g.lastStep = time.Now()
			}
			if sunk && g.userBoard.AllShipsSunk() {
				g.gameOver = true
				g.winner = "Player"
			}

		}
	}
}

func (g *game) step() {
	if g.gameOver {
		return
	}
	time.Sleep(500 * time.Millisecond) // Add a small delay for the AI's turn
	x, y := application.TakeTurn(g.aiViewBoard)
	hit, sunk, shipType, err := g.aiSolutionBoard.Attack(x, y)

	if err != nil {
		fmt.Println("AI error: ", err)
		g.isPlayerTurn = true
		return
	}
	if hit {
		g.aiViewBoard.SetCoordinate(x, y, application.Hit)
	} else {
		g.aiViewBoard.SetCoordinate(x, y, application.Miss)
	}

	if sunk {
		g.aiViewBoard.RecordSunkShip(shipType)
		g.aiSolutionBoard.RecordSunkShip(shipType)
		g.aiViewBoard.CopyHitValues(g.aiSolutionBoard)

		for ints, element := range g.aiSolutionBoard.HitShipAt() {
			if element == shipType {
				g.aiViewBoard.SetCoordinate(ints[0], ints[1], application.SUNK)
			}
		}
		if g.aiSolutionBoard.AllShipsSunk() {
			g.gameOver = true
			g.winner = "AI"
		}
	}

	// After AI move, return control to player
	if !hit {
		g.isPlayerTurn = true
	}
}

func applyAlpha(c color.Color, alpha uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
}

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
		cellSize:     80,
		stepEvery:    time.Millisecond * 100, // kept for consistency; not used yet for turn timing
		rows:         cfg.BATTLESHIPHEIGHT,
		cols:         cfg.BATTLESHIPWIDTH,
		userBoard:    application.NewBattleshipBoard(cfg.BATTLESHIPWIDTH, cfg.BATTLESHIPHEIGHT),
		aiBoard:      application.NewBattleshipBoard(cfg.BATTLESHIPWIDTH, cfg.BATTLESHIPHEIGHT),
		isPlayerTurn: true,
		ai:           application.NewSimpleAI(cfg.BATTLESHIPWIDTH, cfg.BATTLESHIPHEIGHT),
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	g.aiBoard.SeedBoard()
	g.aiBoard.PrintBoard()
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
	rows, cols    int
	cellSize      int
	stepEvery     time.Duration
	lastStep      time.Time
	userBoard     application.BattleshipBoard
	aiBoard       application.BattleshipBoard
	isPlayerTurn bool
	ai            application.AI
}

func (g *game) Update() error {
	if g.isPlayerTurn {
		g.handleClick()
	} else {
		g.step()
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// Clear
	screen.Fill(color.RGBA{A: 255})

	// Colors
	grid := color.RGBA{R: 220, G: 220, B: 220, A: 255}
	userColor := color.RGBA{R: 200, G: 230, B: 255, A: 255}
	shotColor := color.RGBA{R: 255, G: 80, B: 80, A: 255}
	missColor := color.RGBA{R: 120, G: 120, B: 120, A: 255}

	cs := g.cellSize
	boardW := g.cols * cs
	gap := 20

	// Draw user board (left)
	drawGrid(screen, 0, 0, g.cols, g.rows, cs, grid)
	// Fill any AI shots on user board (placeholder: none yet)
	for x := 0; x < g.cols; x++ {
		for y := 0; y < g.rows; y++ {
			chosenColor := userColor
			if g.userBoard.Coordinate(x, y) == application.Miss {
				chosenColor = missColor
			}
			if g.userBoard.Coordinate(x, y) == application.Hit {
				chosenColor = shotColor
			}

			xPix := x*cs + 1
			yPix := y*cs + 1
			vector.DrawFilledRect(screen, float32(xPix), float32(yPix), float32(cs-2), float32(cs-2), chosenColor, false)
		}
	}

	// Draw AI board (bottom)
	offsetY := boardW + gap
	drawGrid(screen, 0, offsetY, g.cols, g.rows, cs, grid)
	// Draw user shots on AI board
	for x := 0; x < g.cols; x++ {
		for y := 0; y < g.rows; y++ {
			chosenColor := userColor
			if g.aiBoard.Coordinate(x, y) == application.Miss {
				chosenColor = missColor
			}
			if g.aiBoard.Coordinate(x, y) == application.Hit {
				chosenColor = shotColor
			}

			xPix := x*cs + 1
			yPix := offsetY + y*cs + 1
			vector.DrawFilledRect(screen, float32(xPix), float32(yPix), float32(cs-2), float32(cs-2), chosenColor, false)
		}
	}

	// UI text
	msg := fmt.Sprintf("User board (top) | AI board (bottom)   Cells: %dx%d  CellSize: %d", g.cols, g.rows, g.cellSize)
	op := &text.DrawOptions{}
	op.GeoM.Translate(10, float64(g.rows*g.cellSize*2+gap-28))
	op.ColorScale.ScaleWithColor(color.RGBA{0, 0, 0, 255})
	text.Draw(screen, msg, &text.GoTextFace{Source: mplusFaceSource, Size: 18}, op)
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
			key := [2]int{gridX, gridY}
			// Toggle a placeholder shot state; hit simulation toggles true/false
			hit, sunk, err := g.aiBoard.Attack(key[0], key[1])

			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				g.isPlayerTurn = false
			}
			if hit {
				fmt.Println("Hit!")
			}
			if sunk {
				fmt.Println("Sunk!")
			}

		}
	}
}

func (g *game) step() {
	time.Sleep(500 * time.Millisecond) // Add a small delay for the AI's turn
	x, y := g.ai.TakeTurn(g.userBoard)
	_, _, err := g.userBoard.Attack(x, y)
	if err != nil {
		fmt.Println("AI error: ", err)
	}
	g.isPlayerTurn = true
}

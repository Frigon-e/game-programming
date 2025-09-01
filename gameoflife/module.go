package gameoflife

import (
	"SideProjectGames/gameoflife/internal/ddd"
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

// Run launches an Ebiten window to visualize Conway's Game of Life using the provided config.
// It replaces the previous CLI printing loop with a graphical, interactive loop.
func Run(cfg config.AppConfig) error {
	g := &game{
		cellSize:  10,
		stepEvery: time.Millisecond * 100,
		read:      ddd.NewGOLBoard(cfg.GOLWIDTH, cfg.GOLHEIGHT),
		write:     ddd.NewGOLBoard(cfg.GOLWIDTH, cfg.GOLHEIGHT),
	}
	g.read.SeedBoard()

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	mplusFaceSource = s

	// Size in pixels
	w := cfg.GOLWIDTH * g.cellSize
	h := cfg.GOLHEIGHT * g.cellSize
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Conway's Game of Life")

	return ebiten.RunGame(g)
}

type skippableItems struct {
	row int16
	col int16
}

type game struct {
	read      ddd.GolBoard
	write     ddd.GolBoard
	skipCord  []skippableItems
	cellSize  int
	stepEvery time.Duration
	lastStep  time.Time
}

func (g *game) Update() error {
	// Step the simulation at fixed intervals
	g.handleClick()
	if time.Since(g.lastStep) >= g.stepEvery {
		g.step()
		g.lastStep = time.Now()
		g.wipeSkippable()
	}
	return nil
}

func (g *game) addSkippable(item skippableItems) {
	g.skipCord = append(g.skipCord, item)
}

func (g *game) wipeSkippable() {
	g.skipCord = []skippableItems{}
}

func (g *game) Draw(screen *ebiten.Image) {
	// Clear
	screen.Fill(color.RGBA{A: 255})

	// Draw alive cells as white rectangles
	cs := g.cellSize
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for r := 0; r < g.read.Rows(); r++ {
		for c := 0; c < g.read.Cols(); c++ {
			if g.read.Coordinate(r, c) {
				x := c * cs
				y := r * cs

				vector.DrawFilledRect(screen, float32(x), float32(y), float32(cs-1), float32(cs-1), white, false)
			}
		}
	}

	//msg := strconv.FormatInt(int64(g.stepEvery), 10)
	msg := fmt.Sprintf("Step Time: %v", g.stepEvery)

	textSize, _ := text.Measure(msg, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}, 24)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(g.read.Rows()*g.cellSize)-textSize, 10)
	op.ColorScale.ScaleWithColor(color.RGBA{255, 0, 0, 255})
	text.Draw(screen, msg, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   24,
	}, op)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *game) handleClick() {
	mouseX, mouseY := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		gridX := mouseX / g.cellSize
		gridY := mouseY / g.cellSize

		if gridX >= 0 && gridX < g.read.Cols() && gridY >= 0 && gridY < g.read.Rows() {
			g.addSkippable(skippableItems{int16(gridY), int16(gridX)})
			current := g.read.Coordinate(gridY, gridX)
			val := !current
			g.write.SetCoordinate(gridY, gridX, val)
			g.read.SetCoordinate(gridY, gridX, val)

		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if g.stepEvery <= time.Millisecond*1000 {
			g.stepEvery += time.Millisecond * 10
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if g.stepEvery >= time.Millisecond*10 {
			g.stepEvery -= time.Millisecond * 10
		}
	}
}

func (g *game) step() {
	// Apply Conway rules from read -> write, then copy back
	for r := 0; r < g.read.Rows(); r++ {
		for c := 0; c < g.read.Cols(); c++ {
			skipStep := false
			for _, items := range g.skipCord {
				if items.row == int16(r) && items.col == int16(c) {
					skipStep = true
					break
				}
			}

			if skipStep {
				continue
			}

			alive := g.read.Coordinate(r, c)
			neighbors := g.read.CountSurroundingLive(r, c)

			newVal := alive
			if alive && neighbors < 2 {
				newVal = false
			}
			if alive && (neighbors == 2 || neighbors == 3) {
				newVal = true
			}
			if alive && neighbors > 3 {
				newVal = false
			}
			if !alive && neighbors == 3 {
				newVal = true
			}
			g.write.SetCoordinate(r, c, newVal)
		}
	}
	g.read.CopyBoard(g.write.FlatSlice())
}

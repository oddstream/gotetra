// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ScreenWidth and ScreenHeight are exported to main.go for ebiten.SetWindowSize()
const (
	ScreenWidth  = 700
	ScreenHeight = 1000
	boardWidth   = 6
	boardHeight  = 9
	NORTH        = 1
	EAST         = 2
	SOUTH        = 4
	WEST         = 8
	MASK         = 15
)

// Game represents a game state.
type Game struct {
	grid      *Grid
	gridImage *ebiten.Image
}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{}
	var err error
	g.grid, err = NewGrid(boardWidth, boardHeight)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	if err := g.grid.Update(); err != nil {
		return err
	}
	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// screen.Fill(backgroundColor)

	if g.gridImage == nil {
		w, h := g.grid.Size()
		// println("gridImage", w, h)
		g.gridImage = ebiten.NewImage(w, h)
	}
	// center gridImage in the screen
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.gridImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.gridImage, op)

	g.grid.Draw(g.gridImage)
}

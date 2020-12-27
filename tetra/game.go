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
	ScreenWidth  = 640
	ScreenHeight = 480
	boardWidth   = 6
	boardHeight  = 8
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
	if g.gridImage == nil {
		w, h := g.grid.Size()
		g.gridImage = ebiten.NewImage(w, h)
	}
	screen.Fill(backgroundColor)
	g.grid.Draw(g.gridImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.gridImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.gridImage, op)
}

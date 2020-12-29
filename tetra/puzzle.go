// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ScreenWidth and ScreenHeight are exported to main.go for ebiten.SetWindowSize()
const (
	boardWidth  = 6
	boardHeight = 9
	NORTH       = 0b0001
	EAST        = 0b0010
	SOUTH       = 0b0100
	WEST        = 0b1000
	MASK        = 0b1111
)

// Puzzle represents a game state.
type Puzzle struct {
	grid      *Grid
	gridImage *ebiten.Image
}

// Init initializes a Level object that was created by the caller
func (p *Puzzle) Init() {
	var err error
	p.grid, err = NewGrid(boardWidth, boardHeight)
	if err != nil {
		log.Fatal(err)
	}
}

// Layout implements ebiten.Game's Layout.
func (p *Puzzle) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (p *Puzzle) Update() error {
	if err := p.grid.Update(); err != nil {
		return err
	}
	return nil
}

// Draw draws the current game to the given screen.
func (p *Puzzle) Draw(screen *ebiten.Image) {
	// screen.Fill(backgroundColor)

	if p.gridImage == nil {
		w, h := p.grid.Size()
		// println("gridImage", w, h)
		p.gridImage = ebiten.NewImage(w, h)
	}
	// center gridImage in the screen
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := p.gridImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(p.gridImage, op)

	p.grid.Draw(p.gridImage)
}

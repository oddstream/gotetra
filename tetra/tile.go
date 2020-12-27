// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	tilesheetImage *ebiten.Image
)

func init() {
	var err error
	tilesheetImage, _, err = ebitenutil.NewImageFromFile("assets/tilesheet9x9x100.png")
	if err != nil {
		log.Fatal(err)
	}
}

// Tile object describes a tile
type Tile struct {
	X, Y       int
	N, E, S, W *Tile
}

// NewTile creates a new Tile object and returns a pointer to it
func NewTile(x, y int) *Tile {
	t := &Tile{X: x, Y: y}
	return t
}

// Draw handles rendering of Tile object
func (t *Tile) Draw() error {
	return nil
}

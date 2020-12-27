// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Grid is an object representing the grid of tiles
type Grid struct {
	width  int
	height int
	tiles  []Tile
}

func (g *Grid) findTile(x, y int) *Tile {
	for i := 0; i < len(g.tiles); i++ {
		t := &g.tiles[i]
		if t.X == x && t.Y == y {
			return t
		}
	}
	return nil
}

// NewGrid create a Grid object
func NewGrid(w, h int) (*Grid, error) {
	g := &Grid{width: w, height: h}
	g.tiles = make([]Tile, w*h)
	for i := 0; i < len(g.tiles); i++ {
		t := &g.tiles[i]
		x := i % w
		y := i / w
		t.X = x
		t.Y = y
		// g.tiles[i] = NewTile(x, y)
	}
	// for i, t := range g.tiles {
	// 	println(i, t.X, t.Y)
	// }
	for i := 0; i < len(g.tiles); i++ {
		t := &g.tiles[i]
		x := i % w
		y := i / w
		t.N = g.findTile(x, y-1)
		t.E = g.findTile(x+1, y)
		t.S = g.findTile(x, y+1)
		t.W = g.findTile(x-1, y)
	}
	// for i, t := range g.tiles {
	// 	println(i, t.X, t.Y, t.N, t.E, t.S, t.W)
	// }
	t0 := g.findTile(0, 0)
	if t0 == nil {
		log.Fatal("cannot find t0")
	}
	t10 := g.findTile(1, 0)
	if t10 == nil {
		log.Fatal("cannot find t10")
	}
	if t0.E != t10 {
		log.Fatal("t0 t10 not linked")
	}
	if t10.W != t0 {
		log.Fatal("t10 t0 not linked")
	}
	return g, nil
}

// Size returns the size of the grid in pixels
func (g *Grid) Size() (int, int) {
	return g.width * 100, g.height * 100
}

// Update the board state (transitions, user input)
func (g *Grid) Update() error {
	return nil
}

// Draw renders the grid into the gridImage
func (g *Grid) Draw(gridImage *ebiten.Image) {
	// display the background
	gridImage.Fill(backgroundColor)
	// then tell each tile to draw itself
	for i := 0; i < len(g.tiles); i++ {
		t := &g.tiles[i]
		t.Draw()
	}
}

// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Grid is an object representing the grid of tiles
type Grid struct {
	width  int
	height int
	tiles  []*Tile // a slice (not array!) of pointers to Tile objects
}

func (g *Grid) findTile(x, y int) *Tile {
	for _, t := range g.tiles {
		if t.X == x && t.Y == y {
			return t
		}
	}
	return nil
}

// NewGrid create a Grid object
func NewGrid(w, h int) (*Grid, error) {
	g := &Grid{width: w, height: h, tiles: make([]*Tile, w*h)}
	for i := range g.tiles {
		g.tiles[i] = NewTile(i%w, i/w)
	}
	// for i, t := range g.tiles {
	// 	println(i, t.X, t.Y)
	// }
	for _, t := range g.tiles {
		x := t.X
		y := t.Y
		t.N = g.findTile(x, y-1)
		t.E = g.findTile(x+1, y)
		t.S = g.findTile(x, y+1)
		t.W = g.findTile(x-1, y)
	}

	// for i, t := range g.tiles {
	// 	println(i, t.X, t.Y, t.N, t.E, t.S, t.W)
	// }
	t00 := g.findTile(0, 0)
	if t00 == nil {
		log.Fatal("cannot find t00")
	}
	t10 := g.findTile(1, 0)
	if t10 == nil {
		log.Fatal("cannot find t10")
	}
	if t00.E != t10 {
		log.Fatal("t00 t10 not linked")
	}
	if t10.W != t00 {
		log.Fatal("t10 t0 not linked")
	}
	t01 := g.findTile(0, 1)
	if t01 == nil {
		log.Fatal("cannot find t01")
	}
	if t00.S != t01 {
		log.Fatal("t00 t01 not linked")
	}
	if t00.S.N != t00 {
		log.Fatal("t00.S.N not linked")
	}
	for _, t := range g.tiles {
		t.PlaceCoin()
	}
	for _, t := range g.tiles {
		t.Jumble()
		t.SetImage()
	}

	return g, nil
}

// Size returns the size of the grid in pixels
func (g *Grid) Size() (int, int) {
	return g.width * TileWidth, g.height * TileHeight
}

// FindTileAt finds the tile under the mouse click or touch
func (g *Grid) FindTileAt(x, y int) *Tile {
	for _, t := range g.tiles {
		x0, y0, x1, y1 := t.Rect()
		if x > x0 && x < x1 && y > y0 && y < y1 {
			return t
		}
	}
	return nil
}

// Update the board state (transitions, user input)
func (g *Grid) Update() error {
	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		originX := (ScreenWidth - (g.width * TileWidth)) / 2
		originY := (ScreenHeight - (g.height * TileHeight)) / 2
		// println(originX, originY, x, y)
		tile := g.FindTileAt(x-originX, y-originY)
		if tile != nil {
			tile.Rotate()
		}
	}

	for _, t := range g.tiles {
		t.Update()
	}
	return nil
}

// Draw renders the grid into the gridImage
func (g *Grid) Draw(gridImage *ebiten.Image) {
	// display the background
	gridImage.Fill(backgroundColor)
	// then tell each tile to draw itself on the gridImage
	for _, t := range g.tiles {
		t.Draw(gridImage)
	}
}

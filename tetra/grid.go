// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

var acmeLargeFont font.Face

func init() {
	bytes, err := ioutil.ReadFile("assets/Acme-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	// https://pkg.go.dev/golang.org/x/image@v0.0.0-20201208152932-35266b937fa6/font
	tt, err := truetype.Parse(bytes)
	if err != nil {
		log.Fatal(err)
	}
	acmeLargeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    144,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// Grid is an object representing the grid of tiles
type Grid struct {
	width  int
	height int
	tiles  []*Tile       // a slice (not array!) of pointers to Tile objects
	colors []*color.RGBA // a slice of pointers to colors for the tiles, one color per section
	level  int
}

func (g *Grid) findTile(x, y int) *Tile {
	// for _, t := range g.tiles {
	// 	if t.X == x && t.Y == y {
	// 		return t
	// 	}
	// }
	// return nil
	if x < 0 || x >= g.width {
		return nil
	}
	if y < 0 || y >= g.height {
		return nil
	}
	i := x + (y * g.width)
	if i < 0 || i > len(g.tiles) {
		log.Fatal("findTile index out of bounds")
	}
	return g.tiles[i]
}

// NewGrid create a Grid object
func NewGrid(w, h int) (*Grid, error) {
	g := &Grid{width: w, height: h, tiles: make([]*Tile, w*h)}
	for i := range g.tiles {
		g.tiles[i] = NewTile(g, i%w, i/w)
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

	/*
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
	*/

	g.colors = []*color.RGBA{
		&colorNavy,
		&colorBlue,
		&colorCornflowerBlue,
		&colorLightSkyBlue,
	} // golang gotcha no newline after last literal, must be comma or closing brace

	// TODO load level from saved state
	g.level = 0
	g.NextLevel()

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

func (g *Grid) findUnsectionedTile() *Tile {
	for _, t := range g.tiles {
		if t.coins != 0 && t.color == nil {
			return t
		}
	}
	return nil
}

// ColorTiles sets the color and section for every tile
func (g *Grid) ColorTiles() {
	nextColor := 0
	nextSection := 0
	tile := g.findUnsectionedTile()
	if tile == nil {
		panic("no first unsection tile")
	}
	for tile != nil {
		tile.ColorConnected(g.colors[nextColor], nextSection)
		nextColor++
		if nextColor >= len(g.colors) {
			nextColor = 0
		}
		nextSection++
		tile = g.findUnsectionedTile()
	}
}

// IsSectionComplete returns true if all the tiles in a section are aligned
func (g *Grid) IsSectionComplete(section int) bool {
	for _, t := range g.tiles {
		if t.section == section {
			if !t.IsCompleteSection(section) {
				return false
			}
		}
	}
	return true
}

// TriggerScaleDown starts the tile disappear process
func (g *Grid) TriggerScaleDown(section int) {
	for _, t := range g.tiles {
		if t.section == section {
			t.targScale = 0.1
		}
	}
}

// FilterSection applies a Tile function to all tiles in the section
func (g *Grid) FilterSection(f func(*Tile), section int) {
	for _, t := range g.tiles {
		if t.section == section {
			f(t)
		}
	}
}

// IsComplete returns true if all the tiles are aligned
func (g *Grid) IsComplete() bool {
	for _, t := range g.tiles {
		if !t.IsComplete() {
			return false
		}
	}
	return true
}

// NextLevel resets game data and moves the puzzle to the next level
func (g *Grid) NextLevel() {
	for _, t := range g.tiles {
		t.Reset()
	}

	g.level++
	rand.Seed(int64(g.level))

	for _, t := range g.tiles {
		t.PlaceCoin()
	}
	g.ColorTiles()
	for _, t := range g.tiles {
		t.Jumble()
		t.SetImage()
	}
}

// Update the board state (transitions, user input)
func (g *Grid) Update() error {
	// TODO move input up to puzzle so we can reset at that level (where level is)
	// or just move level down here
	// if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.IsComplete() {
			g.NextLevel()
		} else {
			x, y := ebiten.CursorPosition()
			originX := (ScreenWidth - (g.width * TileWidth)) / 2
			originY := (ScreenHeight - (g.height * TileHeight)) / 2
			// println(originX, originY, x, y)
			tile := g.FindTileAt(x-originX, y-originY)
			if tile != nil {
				tile.Rotate()
			}
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

	str := fmt.Sprint(g.level)
	bound, _ := font.BoundString(acmeLargeFont, str)
	w := (bound.Max.X - bound.Min.X).Ceil()
	h := (bound.Max.Y - bound.Min.Y).Ceil()
	x, y := g.Size()
	x = (x / 2) - (w / 2)
	y = (y / 2) + (h / 2)
	colorTransBlack := color.RGBA{R: 0, G: 0, B: 0, A: 0x10}
	text.Draw(gridImage, str, acmeLargeFont, x, y, colorTransBlack)

	for _, t := range g.tiles {
		t.Draw(gridImage)
	}
}

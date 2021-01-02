// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Grid is an object representing the grid of tiles
type Grid struct {
	mode    string // "bubblewrap" | "puzzle"
	width   int
	height  int
	tiles   []*Tile // a slice (not array!) of pointers to Tile objects
	palette Palette
	colors  []*color.RGBA // a slice of pointers to colors for the tiles, one color per section
	ud      *UserData
	spores  []*Spore
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
func NewGrid(m string, w, h int) (*Grid, error) {
	g := &Grid{mode: m, width: w, height: h, tiles: make([]*Tile, w*h)}
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

	g.ud = NewUserData()
	// NextLevel will bump the UserData.Level, which isn't what we want, so
	g.ud.Level--
	g.NextLevel()

	g.spores = make([]*Spore, 0, 32)

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
		if t.coins != 0 && t.color == BasicColors["Black"] {
			return t
		}
	}
	return nil
}

// ColorTiles sets the color and section for every tile
func (g *Grid) ColorTiles() {
	switch g.mode {
	case "bubblewrap":
		nextColor := 0
		nextSection := 0
		tile := g.findUnsectionedTile()
		if tile == nil {
			panic("no first unsection tile")
		}
		for tile != nil {
			tile.ColorConnected(g.palette[nextColor], nextSection)
			nextColor++
			if nextColor == len(g.palette) {
				nextColor = 0
			}
			nextSection++
			tile = g.findUnsectionedTile()
		}
	case "puzzle":
		for _, t := range g.tiles {
			t.section = 0 // any number will do
			{
				n := g.ud.Level % len(g.palette)
				colName := g.palette[n]
				t.color = ExtendedColors[colName]
			}
		}
	default:
		log.Fatal("unknown mode", g.mode)
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
// func (g *Grid) TriggerScaleDown(section int) {
// 	for _, t := range g.tiles {
// 		if t.section == section {
// 			t.targScale = 0.1
// 		}
// 	}
// }

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

	g.ud.Level++
	g.ud.Save()
	rand.Seed(int64(g.ud.Level))

	for _, t := range g.tiles {
		t.PlaceCoin()
	}
	g.palette = Palettes[rand.Int()%len(Palettes)]
	g.ColorTiles()
	for _, t := range g.tiles {
		t.Jumble()
		t.SetImage()
	}
}

// AddSpore to map of spores
func (g *Grid) AddSpore(x, y int, img *ebiten.Image, deg int, col color.RGBA) {
	x *= TileWidth
	x += TileWidth / 2
	y *= TileWidth
	y += TileWidth / 2
	sp := NewSpore(x, y, img, float64(deg), col)
	g.spores = append(g.spores, sp)
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

	for _, sp := range g.spores {
		sp.Update()
	}

	// if len(g.spores) > 16 {
	// 	g.spores = g.spores[1:]
	// 	println("trimmed spores")
	// }
	{
		newSpores := make([]*Spore, 0, len(g.spores))
		for _, sp := range g.spores {
			if sp.IsVisible() {
				newSpores = append(newSpores, sp)
			}
		}
		g.spores = newSpores
	}

	return nil
}

// Draw renders the grid into the gridImage
func (g *Grid) Draw(gridImage *ebiten.Image) {
	// display the background
	gridImage.Fill(backgroundColor)

	str := fmt.Sprint(g.ud.Level)
	bound, _ := font.BoundString(Acme.huge, str)
	w := (bound.Max.X - bound.Min.X).Ceil()
	h := (bound.Max.Y - bound.Min.Y).Ceil()
	x, y := g.Size()
	x = (x / 2) - (w / 2)
	y = (y / 2) + (h / 2)
	colorTransBlack := color.RGBA{R: 0, G: 0, B: 0, A: 0x10}
	text.Draw(gridImage, str, Acme.huge, x, y, colorTransBlack)

	for _, t := range g.tiles {
		t.Draw(gridImage)
	}

	for _, sp := range g.spores {
		sp.Draw(gridImage)
	}
	ebitenutil.DebugPrint(gridImage, fmt.Sprintf("%d spores", len(g.spores)))
}

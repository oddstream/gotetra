// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// TilesAcross and TilesDown are package-level variables so they can be seen by Tile
var (
	TilesAcross int
	TilesDown   int
	LeftMargin  int
	TopMargin   int
	TileSize    int
)

// Grid is an object representing the grid of tiles
type Grid struct {
	mode            string  // "bubblewrap" | "puzzle"
	tiles           []*Tile // a slice (not array!) of pointers to Tile objects
	palette         Palette
	colors          []*color.RGBA // a slice of pointers to colors for the tiles, one color per section
	ud              *UserData
	spores          []*Spore
	colorBackground color.RGBA
}

func (g *Grid) findTile(x, y int) *Tile {
	// for _, t := range g.tiles {
	// 	if t.X == x && t.Y == y {
	// 		return t
	// 	}
	// }
	// return nil
	if x < 0 || x >= TilesAcross {
		return nil
	}
	if y < 0 || y >= TilesDown {
		return nil
	}
	i := x + (y * TilesAcross)
	if i < 0 || i > len(g.tiles) {
		log.Fatal("findTile index out of bounds")
	}
	return g.tiles[i]
}

// NewGrid create a Grid object
func NewGrid(m string, w, h int) *Grid {

	if w == 0 || h == 0 {
		TileSize = 100
		TilesAcross, TilesDown = ScreenWidth/TileSize, ScreenHeight/TileSize
	} else {
		possibleW := ScreenWidth / (w + 1) // add 1 to create margin for endcaps
		possibleH := ScreenHeight / (h + 1)
		// golang gotcha there isn't a vanilla math.MinInt()
		if possibleW < possibleH {
			TileSize = possibleW
		} else {
			TileSize = possibleH
		}
		println("TileSize", TileSize)
		TilesAcross, TilesDown = w, h
	}
	LeftMargin = (ScreenWidth - (TilesAcross * TileSize)) / 2
	TopMargin = (ScreenHeight - (TilesDown * TileSize)) / 2

	// now we know the Size Of Things, tell Tile to load it's static stuff
	initTileImages()

	g := &Grid{mode: m, tiles: make([]*Tile, TilesAcross*TilesDown)}
	for i := range g.tiles {
		g.tiles[i] = NewTile(g, i%TilesAcross, i/TilesAcross)
	}

	// link the tiles together to avoid all that tedious 2d array stuff
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
	g.ud.Level-- // TODO this is ugly, maybe .CompletedLevel?
	g.NextLevel()

	g.spores = make([]*Spore, 0, TilesAcross*TilesDown)

	return g
}

// Size returns the size of the grid in pixels
func (g *Grid) Size() (int, int) {
	return TilesAcross * TileSize, TilesDown * TileSize
}

// FindTileAt finds the tile under the mouse click or touch
func (g *Grid) FindTileAt(pt image.Point) *Tile {
	for _, t := range g.tiles {
		if InRect(pt, t.Rect) {
			return t
		}
	}
	return nil
}

func (g *Grid) findUnsectionedTile() *Tile {
	for _, t := range g.tiles {
		if t.coins != 0 && t.color == colorUnsectioned {
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
	g.colorBackground = CalcBackgroundColor(g.palette)
	g.ColorTiles()
	for _, t := range g.tiles {
		t.Jumble()
		t.SetImage()
	}
}

// AddSpore to map of spores
func (g *Grid) AddSpore(x, y int, img *ebiten.Image, deg int, col color.RGBA) {
	// convert X,Y into screen coords of tile center
	xScreen := LeftMargin + (x * TileSize) + (TileSize / 2)
	yScreen := TopMargin + (y * TileSize) + (TileSize / 2)
	sp := NewSpore(float64(xScreen), float64(yScreen), img, float64(deg), col)
	g.spores = append(g.spores, sp)
}

// Layout implements ebiten.Game's Layout.
func (g *Grid) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update the board state (transitions, user input)
func (g *Grid) Update(i *Input) error {

	i.Update()

	if i.pt.X != 0 && i.pt.Y != 0 {
		if g.IsComplete() {
			g.NextLevel()
		} else {
			// could treat the Tiles as Widgets
			// implement Tile.Pushed(), Tile.Action()
			// would mean asking each tile during Tile.Update()
			// or creating an object that links Input with a list of widgets
			// Grid has Widget[] instead of []Tile
			tile := g.FindTileAt(i.pt)
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
func (g *Grid) Draw(screen *ebiten.Image) {

	screen.Fill(g.colorBackground)

	{
		str := fmt.Sprint(g.ud.Level)
		bound, _ := font.BoundString(Acme.huge, str)
		w := (bound.Max.X - bound.Min.X).Ceil()
		h := (bound.Max.Y - bound.Min.Y).Ceil()
		x, y := g.Size()
		x = (x / 2) - (w / 2)
		y = (y / 2) + (h / 2)
		x += LeftMargin
		y += TopMargin
		colorTransBlack := color.RGBA{R: 0, G: 0, B: 0, A: 0x10}
		text.Draw(screen, str, Acme.huge, x, y, colorTransBlack)
	}

	for _, t := range g.tiles {
		t.Draw(screen)
	}

	for _, sp := range g.spores {
		sp.Draw(screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%d,%d grid of size %d, %d spores", TilesAcross, TilesDown, TileSize, len(g.spores)))

}

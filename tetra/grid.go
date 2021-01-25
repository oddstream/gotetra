// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	frags           []*Frag
	colorBackground color.RGBA
	stroke          *Stroke
}

func (g *Grid) findTile(x, y int) *Tile {

	// we re-order the tiles when dragging, to put the dragged tile at the top of the z-order
	// so we can't use i := x + (y * TilesAcross) to find index of tile in slice
	// except if this func is just used after tiles have been created

	for _, t := range g.tiles {
		if t.X == x && t.Y == y {
			return t
		}
	}
	return nil

	// if x < 0 || x >= TilesAcross {
	// 	return nil
	// }
	// if y < 0 || y >= TilesDown {
	// 	return nil
	// }
	// i := x + (y * TilesAcross)
	// if i < 0 || i > len(g.tiles) {
	// 	log.Fatal("findTile index out of bounds")
	// }
	// return g.tiles[i]
}

// NewGrid create a Grid object
func NewGrid(m string, w, h int) *Grid {

	var screenWidth, screenHeight int

	if runtime.GOARCH == "wasm" {
		screenWidth, screenHeight = WindowWidth, WindowHeight
	} else {
		screenWidth, screenHeight = ebiten.WindowSize()
	}

	if w == 0 || h == 0 {
		TileSize = 100
		TilesAcross, TilesDown = screenWidth/TileSize, screenHeight/TileSize
	} else {
		possibleW := screenWidth / (w + 1) // add 1 to create margin for endcaps
		possibleW /= 20
		possibleW *= 20
		possibleH := screenHeight / (h + 1)
		possibleH /= 20
		possibleH *= 20
		// golang gotcha there isn't a vanilla math.MinInt()
		if possibleW < possibleH {
			TileSize = possibleW
		} else {
			TileSize = possibleH
		}
		TilesAcross, TilesDown = w, h
	}
	LeftMargin = (screenWidth - (TilesAcross * TileSize)) / 2
	TopMargin = (screenHeight - (TilesDown * TileSize)) / 2

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
		t.edges[0] = g.findTile(x, y-1)
		t.edges[1] = g.findTile(x+1, y)
		t.edges[2] = g.findTile(x, y+1)
		t.edges[3] = g.findTile(x-1, y)
	}

	g.ud = NewUserData()
	g.CreateNextLevel()

	g.frags = make([]*Frag, 0, TilesAcross*TilesDown)

	return g
}

// Size returns the size of the grid in pixels
func (g *Grid) Size() (int, int) {
	return TilesAcross * TileSize, TilesDown * TileSize
}

// findTileAt finds the tile under the mouse click or touch
func (g *Grid) findTileAt(pt image.Point) *Tile {
	for _, t := range g.tiles {
		if InRect(pt, t.Rect) {
			return t
		}
	}
	return nil
}

func (g *Grid) moveTileToFront(t *Tile) {
	index := -1
	for i, t0 := range g.tiles {
		if t0 == t {
			index = i
			break
		}
	}
	// https://github.com/golang/go/wiki/SliceTricks
	g.tiles = append(g.tiles[:index], g.tiles[index+1:]...)
	g.tiles = append(g.tiles, t)
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
				n := g.ud.CompletedLevels % len(g.palette)
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
			if !t.IsComplete() {
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

// CreateNextLevel resets game data and moves the puzzle to the next level
func (g *Grid) CreateNextLevel() {
	for _, t := range g.tiles {
		t.Reset()
	}

	rand.Seed(int64(g.ud.CompletedLevels))

	// y3 := TilesDown / 3
	// y6 := y3 + y3
	// {
	// 	tp := &TilePath{start: g.findTile(0, y3)}
	// 	tp.Run(EAST)
	// }
	// {
	// 	tp := &TilePath{start: g.findTile(TilesAcross-1, y6)}
	// 	tp.Run(WEST)
	// }
	// x3 := TilesAcross / 3
	// x6 := x3 + x3
	// {
	// 	tp := &TilePath{start: g.findTile(x3, 0)}
	// 	tp.Run(SOUTH)
	// }
	// {
	// 	tp := &TilePath{start: g.findTile(x6, TilesDown-1)}
	// 	tp.Run(NORTH)
	// }

	for _, t := range g.tiles {
		t.PlaceRandomCoins()
	}

	g.palette = Palettes[rand.Int()%len(Palettes)]
	g.colorBackground = CalcBackgroundColor(g.palette)
	g.ColorTiles()
	for _, t := range g.tiles {
		t.Jumble()
		t.SetImage()
	}
}

// AddFrag to map of frags
func (g *Grid) AddFrag(x, y int, img *ebiten.Image, deg int, col color.RGBA) {
	// convert X,Y into screen coords of tile center
	xScreen := LeftMargin + (x * TileSize) + (TileSize / 2)
	yScreen := TopMargin + (y * TileSize) + (TileSize / 2)
	sp := NewFrag(float64(xScreen), float64(yScreen), img, float64(deg), col)
	g.frags = append(g.frags, sp)
}

// Layout implements ebiten.Game's Layout.
func (g *Grid) Layout(outsideWidth, outsideHeight int) (int, int) {
	LeftMargin = (outsideWidth - (TilesAcross * TileSize)) / 2
	TopMargin = (outsideHeight - (TilesDown * TileSize)) / 2
	for _, t := range g.tiles {
		t.Layout()
	}
	return outsideWidth, outsideHeight
}

// Update the board state (transitions, user input)
func (g *Grid) Update() error {

	if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) {
		GSM.Switch(NewSplash())
	}

	if g.stroke == nil {
		_, yoff := ebiten.Wheel()
		if yoff != 0 {
			x, y := ebiten.CursorPosition()
			t := g.findTileAt(image.Point{X: x, Y: y})
			if t != nil {
				if yoff < 0 {
					t.Rotate()
				} else {
					t.Unrotate()
				}
			}
			return nil
		}
	}

	var s *Stroke

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s = NewStroke(&MouseStrokeSource{})
	}
	ts := inpututil.JustPressedTouchIDs()
	if ts != nil && len(ts) == 1 {
		s = NewStroke(&TouchStrokeSource{ts[0]})
	}

	if s != nil {
		t := g.findTileAt(s.Position())
		if t != nil && t.state == TileSettled {
			g.stroke = s
			t.state = TileBeingDragged
			g.stroke.SetDraggingObject(t)
			g.moveTileToFront(t)
		}
	}

	if g.stroke != nil {

		g.stroke.Update()

		{
			t := g.stroke.DraggingObject().(*Tile)
			pt := g.stroke.PositionDiff()
			t.offsetX = float64(pt.X)
			t.offsetY = float64(pt.Y)
		}

		if g.stroke.IsReleased() {

			src := g.stroke.DraggingObject().(*Tile)
			if src == nil {
				panic("no tile being dragged")
			}
			dst := g.findTileAt(g.stroke.Position())

			if dst == nil {
				// being dragged off screen
				// beware: tile will use offsetX,Y
				src.state = TileReturning
			} else if src == dst {
				// treat this like a tap
				src.offsetX, src.offsetY = 0, 0 // snap back home
				src.state = TileSettled         // otherwise Tile won't rotate
				if g.IsComplete() {
					g.ud.CompletedLevels++
					g.ud.Save()
					g.CreateNextLevel()
				} else {
					src.Rotate()
				}
			} else if dst.coins == 0 {
				// TODO refactor
				dst.coins = src.coins
				dst.color = src.color
				dst.section = src.section
				dst.tileImage = tileImageLibrary[dst.coins]
				dst.state = TileSettled

				src.coins = 0
				src.color = colorUnsectioned
				src.section = 0
				src.tileImage = tileImageLibrary[0]
				src.offsetX, src.offsetY = 0, 0
				src.state = TileSettled

				if g.IsSectionComplete(dst.section) {
					g.FilterSection((*Tile).TriggerScaleDown, dst.section)
				}
			} else {
				// beware: tile will use offsetX,Y
				src.state = TileReturning
			}
			g.stroke = nil
		}
		// else the stroke isn't released, so the tile is being dragged
	}

	for _, t := range g.tiles {
		t.Update()
	}

	for _, sp := range g.frags {
		sp.Update()
	}

	{
		newFrags := make([]*Frag, 0, len(g.frags))
		for _, sp := range g.frags {
			if sp.IsVisible() {
				newFrags = append(newFrags, sp)
			}
		}
		g.frags = newFrags
	}

	return nil
}

// Draw renders the grid into the gridImage
func (g *Grid) Draw(screen *ebiten.Image) {

	screen.Fill(g.colorBackground)

	{
		str := fmt.Sprint(g.ud.CompletedLevels + 1)
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

	for _, sp := range g.frags {
		sp.Draw(screen)
	}

	for _, t := range g.tiles {
		t.Draw(screen)
	}

	if DebugMode {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%d,%d grid, tile size %d, %d frags", TilesAcross, TilesDown, TileSize, len(g.frags)))
	}
}

// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TileWidth is the unscaled width of a tile in pixels
const TileWidth int = 100

// TileHeight is the unscaled height of a tile in pixels
const TileHeight int = 100

// NORTH is the bit pattern for the upwards direction
const (
	NORTH = 0b0001
	EAST  = 0b0010
	SOUTH = 0b0100
	WEST  = 0b1000
	MASK  = 0b1111
)

type imageInfo struct{ img, deg int }

var (
	coin2ImageInfoMap = map[uint]imageInfo{
		0:                           {0, 0},
		NORTH:                       {1, 180},
		EAST:                        {1, -90},
		SOUTH:                       {1, 0},
		WEST:                        {1, 90},
		NORTH | SOUTH:               {2, 90},
		EAST | WEST:                 {2, 0},
		NORTH | EAST | WEST | SOUTH: {3, 0},
		NORTH | EAST | SOUTH:        {4, 90},
		EAST | SOUTH | WEST:         {4, 180},
		SOUTH | WEST | NORTH:        {4, -90},
		WEST | NORTH | EAST:         {4, 0},
		NORTH | EAST:                {5, 90},
		NORTH | WEST:                {5, 0},
		SOUTH | EAST:                {5, 180},
		SOUTH | WEST:                {5, -90},
	}

	tileImages map[int]*ebiten.Image
	// mplusNormalFont font.Face
)

func getSubImageAndScaleDown(tilesheetImage *ebiten.Image, rect image.Rectangle) *ebiten.Image {

	// had a spot of bother scaling/rotating the tile image in Draw(), so pre-scale the tile images here
	// extract sub image, scale it, draw it into another image, then draw that constucted image into gridImage

	subImage := tilesheetImage.SubImage(rect).(*ebiten.Image)

	// each subImage is 400x400, but it need to appear to be 300x300 when scaled into a 100x100 tile
	scaledImage := ebiten.NewImage(TileWidth, TileHeight)
	op := &ebiten.DrawImageOptions{}
	scaleX := float64(TileWidth) / 400.0
	scaleY := float64(TileHeight) / 400.0
	op.GeoM.Scale(scaleX, scaleY)

	scaledImage.DrawImage(subImage, op)

	return scaledImage
}

func init() {
	tilesheetImage, _, err := ebitenutil.NewImageFromFile("/home/gilbert/Tetra/assets/tilesheet2.png")
	if err != nil {
		log.Fatal(err)
	}

	/*
		0	0, 0	blank
		1	1, 0	one bit	(short line and circle)
		2	2, 0	two bits (line)
		3	0, 1	four bits (cross)
		4	1, 1	three bits
		5	2, 1	two bits (l-shape)

		map from type (0..5) to image.Rect
	*/

	spriteMapRects := map[int]image.Rectangle{
		0: image.Rect(0, 0, 400, 400),
		1: image.Rect(0, 400, 800, 800),
		2: image.Rect(400, 0, 800, 400),
		3: image.Rect(800, 400, 1200, 800),
		4: image.Rect(400, 400, 800, 800),
		5: image.Rect(800, 0, 1200, 800),
	}

	tileImages = map[int]*ebiten.Image{
		0: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[0]),
		1: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[1]),
		2: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[2]),
		3: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[3]),
		4: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[4]),
		5: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[5]),
	}
}

// Tile object describes a tile
type Tile struct {
	G          *Grid
	X, Y       int
	N, E, S, W *Tile

	tileImage   *ebiten.Image
	currDegrees int
	targDegrees int
	currScale   float64
	targScale   float64
	coins       uint
	section     int
	color       *color.RGBA

	// if currDegrees != targDegrees then tile is lerping/rotating
	// if currScale != targScale then tile is lerping/disappearing

	// rotating tile does not receive input

	// don't need hammingWeight because graphics not created dynamically
}

// NewTile creates a new Tile object and returns a pointer to it
func NewTile(g *Grid, x, y int) *Tile {
	t := &Tile{G: g, X: x, Y: y}
	return t
}

// Reset prepares a Tile for a new level by resetting just gameplay data, not structural data
func (t *Tile) Reset() {
	t.coins = 0
	t.section = 0
	t.color = nil

	t.SetImage() // reset to a blank tile image
}

// PlaceCoin sets up a random config for this tile
func (t *Tile) PlaceCoin() {
	bits := [4]uint{NORTH, EAST, SOUTH, WEST}
	opps := [4]uint{SOUTH, WEST, NORTH, EAST}
	links := [4]*Tile{t.N, t.E, t.S, t.W}

	// t.coins = 0
	for d := 0; d < 4; d++ {
		if rand.Float64() < 0.2 {
			if links[d] != nil {
				t.coins |= bits[d]
				links[d].coins |= opps[d]
			}
		}
	}
}

// ColorConnected assigns color and section to tiles connected (by coinage) to this tile
func (t *Tile) ColorConnected(color *color.RGBA, section int) {
	bits := [4]uint{NORTH, EAST, SOUTH, WEST}
	links := [4]*Tile{t.N, t.E, t.S, t.W}

	if color == nil {
		panic("nil color")
	}

	t.color = color
	t.section = section

	for d := 0; d < 4; d++ {
		if t.coins&bits[d] != 0 {
			tn := links[d]
			if tn != nil && tn.coins != 0 && tn.color == nil {
				tn.ColorConnected(color, section)
			}
		}
	}
}

// SetImage is used when all coins are placed
func (t *Tile) SetImage() {
	info := coin2ImageInfoMap[t.coins]
	t.tileImage = tileImages[info.img]
	t.currDegrees = info.deg
	t.targDegrees = info.deg
	t.currScale = 1.0
	t.targScale = 1.0
}

// Rect gives the x,y coords of the tile's top left and bottom right corners, in screen coordinates
func (t *Tile) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = t.X * TileWidth
	y0 = t.Y * TileHeight
	x1 = x0 + TileWidth
	y1 = y0 + TileHeight
	return // using named return parameters
}

func shiftBits(num uint) uint {
	if num&0b1000 == 0b1000 {
		num = num << 1
		num = num & 0b1111
		num = num | 1
	} else {
		num = num << 1
	}
	return num
}

func unshiftBits(num uint) uint {
	if num&1 == 1 {
		num = num >> 1
		num = num | 0b1000
	} else {
		num = num >> 1
	}
	return num
}

// Jumble shifts the bits in the tile a random number of times
func (t *Tile) Jumble() {
	r := rand.Float64()
	if r < 0.25 {
		t.coins = shiftBits(t.coins)
	} else if r < 0.5 {
		t.coins = unshiftBits(t.coins)
	} else if r < 0.75 {
		t.coins = shiftBits(t.coins)
		t.coins = shiftBits(t.coins)
	}
	// TODO remove debugging jumble before release
	// if t.coins == NORTH || t.coins == SOUTH {
	// 	t.coins = unshiftBits(t.coins)
	// }
}

// Rotate shifts the tile 90 degrees clockwise
func (t *Tile) Rotate() {
	if 0 == t.coins {
		return
	}
	if t.currDegrees != t.targDegrees {
		return
	}
	t.coins = shiftBits(t.coins)
	t.targDegrees = t.currDegrees + 90
	if t.targDegrees >= 360 {
		t.targDegrees = 0
	}
	// println("rotate", t.X, t.Y, t.coins)
}

// IsComplete returns true if the tile aligns properly with it's neighbours
func (t *Tile) IsComplete() bool {
	if t.currDegrees != t.targDegrees {
		return false
	}
	if (t.coins & NORTH) == NORTH {
		if (t.N == nil) || ((t.N.coins & SOUTH) == 0) {
			return false
		}
	}
	if (t.coins & EAST) == EAST {
		if (t.E == nil) || ((t.E.coins & WEST) == 0) {
			return false
		}
	}
	if (t.coins & SOUTH) == SOUTH {
		if (t.S == nil) || ((t.S.coins & NORTH) == 0) {
			return false
		}
	}
	if (t.coins & WEST) == WEST {
		if (t.W == nil) || ((t.W.coins & EAST) == 0) {
			return false
		}
	}
	return true
}

// IsCompleteSection returns true if the tile aligns properly with it's neighbours
func (t *Tile) IsCompleteSection(section int) bool {
	// By design, Go does not support optional parameters, default parameter values, or method overloading.
	if t.currDegrees != t.targDegrees {
		return false
	}
	if section != t.section {
		return false
	}
	if (t.coins & NORTH) == NORTH {
		if (t.N == nil) || (t.N.section != section) || ((t.N.coins & SOUTH) == 0) {
			return false
		}
	}
	if (t.coins & EAST) == EAST {
		if (t.E == nil) || (t.E.section != section) || ((t.E.coins & WEST) == 0) {
			return false
		}
	}
	if (t.coins & SOUTH) == SOUTH {
		if (t.S == nil) || (t.S.section != section) || ((t.S.coins & NORTH) == 0) {
			return false
		}
	}
	if (t.coins & WEST) == WEST {
		if (t.W == nil) || (t.W.section != section) || ((t.W.coins & EAST) == 0) {
			return false
		}
	}
	return true
}

// TriggerScaleDown tells this tile to start scaling down
func (t *Tile) TriggerScaleDown() {
	t.targScale = 0.1
}

// Update the tile state (transitions, user input)
func (t *Tile) Update() error {

	if 0 == t.coins {
		return nil
	}

	if t.currDegrees != t.targDegrees {
		t.currDegrees += 5
		if t.currDegrees >= 360 {
			t.currDegrees = 0
		}
		if t.currDegrees == t.targDegrees {
			t.SetImage()
			if t.G.IsSectionComplete(t.section) {
				// t.G.TriggerScaleDown(t.section)
				t.G.FilterSection((*Tile).TriggerScaleDown, t.section)
			}
		}
	}

	if t.targScale < 1.0 {
		t.currScale -= 0.05
		if t.currScale <= t.targScale {
			t.Reset()
		}
	}

	return nil
}

// Draw handles rendering of Tile object
func (t *Tile) Draw(gridImage *ebiten.Image) error {

	// scale, point translation, rotate, object translation

	op := &ebiten.DrawImageOptions{}

	if t.currDegrees != 0 {
		op.GeoM.Translate(-float64(TileWidth)/2.0, -float64(TileHeight)/2.0)
		op.GeoM.Rotate(float64(float64(t.currDegrees) * 3.1415926535 / float64(180)))
		op.GeoM.Translate(float64(TileWidth)/2.0, float64(TileHeight)/2.0)
	}

	// Reset RGB (not Alpha) forcibly
	if t.color != nil {
		op.ColorM.Scale(0, 0, 0, t.currScale)
		// op.ColorM.Scale(0, 0, 0, 1)
		r := float64(t.color.R) / 0xff
		g := float64(t.color.G) / 0xff
		b := float64(t.color.B) / 0xff
		op.ColorM.Translate(r, g, b, 0)
		// op.CompositeMode = ebiten.CompositeModeLighter
	}

	x := float64(t.X * TileWidth)
	y := float64(t.Y * TileHeight)

	/*
		if t.currScale > t.targScale {
			// TODO understand why this works/doesn't work
			// it keeps the tile stable, but the shapes go up to the left
			x += float64(TileWidth/2) * (1.0 - t.currScale)
			y += float64(TileHeight/2) * (1.0 - t.currScale)
			op.GeoM.Scale(t.currScale, t.currScale)
			x -= float64(TileWidth/2) * (1.0 - t.currScale)
			y -= float64(TileHeight/2) * (1.0 - t.currScale)
		}
	*/

	// first move the origin to the center of the tile
	op.GeoM.Translate(-float64(TileWidth/2), -float64(TileHeight/2))
	// scale tile image up to allow endcaps to overhang
	op.GeoM.Scale(400.0/300.0, 400.0/300.0)
	// then move the origin back to top left
	op.GeoM.Translate(float64(TileWidth/2), float64(TileHeight/2))

	op.GeoM.Translate(float64(x), float64(y))

	if t.X%2 == 0 && t.Y%2 == 0 {
		colorTransBlack := color.RGBA{R: 0, G: 0, B: 0, A: 0x10}
		ebitenutil.DrawRect(gridImage, float64(x), float64(y), float64(TileWidth), float64(TileHeight), colorTransBlack)
	}
	gridImage.DrawImage(t.tileImage, op)

	// if t.coins != 0 {
	// 	str := fmt.Sprintf("%04b", t.coins)
	// 	bound, _ := font.BoundString(Acme.small, str)
	// 	w := (bound.Max.X - bound.Min.X).Ceil()
	// 	h := (bound.Max.Y - bound.Min.Y).Ceil()
	// 	x = x + (TileWidth-w)/2
	// 	y = y + (TileHeight-h)/2 + h
	// 	var c color.Color
	// 	if t.IsComplete() {
	// 		c = colorGold
	// 	} else {
	// 		c = colorTeal
	// 	}
	// 	text.Draw(gridImage, str, Acme.small, x, y, c)
	// }

	return nil
}

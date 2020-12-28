// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TileWidth is the unscaled width of a tile in pixels
const TileWidth int = 100

// TileHeight is the unscaled height of a tile in pixels
const TileHeight int = 100

type imageInfo struct{ img, deg int }

var (
	coin2ImageInfoMap = map[uint]imageInfo{
		0:                           {0, 0},
		NORTH:                       {1, 180},
		EAST:                        {1, 270},
		SOUTH:                       {1, 0},
		WEST:                        {1, 90},
		NORTH | SOUTH:               {2, 0},
		EAST | WEST:                 {2, 90},
		NORTH | EAST | WEST | SOUTH: {3, 0},
		NORTH | EAST | SOUTH:        {4, 90},
		EAST | SOUTH | WEST:         {4, 180},
		SOUTH | WEST | NORTH:        {4, 270},
		WEST | NORTH | EAST:         {4, 0},
		NORTH | EAST:                {5, 90},
		NORTH | WEST:                {5, 0},
		SOUTH | EAST:                {5, 180},
		SOUTH | WEST:                {5, 270},
	}

	tileImages map[int]*ebiten.Image
)

func getSubImageAndScaleDown(tilesheetImage *ebiten.Image, rect image.Rectangle) *ebiten.Image {
	subImage := tilesheetImage.SubImage(rect).(*ebiten.Image)

	scaledImage := ebiten.NewImage(TileWidth, TileHeight)
	op := &ebiten.DrawImageOptions{}
	scaleX := float64(TileWidth) / float64(300)
	scaleY := float64(TileHeight) / float64(300)
	op.GeoM.Scale(scaleX, scaleY)
	scaledImage.DrawImage(subImage, op)

	return scaledImage
}

func init() {
	tilesheetImage, _, err := ebitenutil.NewImageFromFile("assets/tilesheet9x9x100.png")
	if err != nil {
		log.Fatal(err)
	}
	const spriteSize int = 300

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
		0: image.Rect(0, 0, 300, 300),
		1: image.Rect(300, 0, 600, 300),
		2: image.Rect(600, 0, 900, 300),
		3: image.Rect(0, 300, 300, 600),
		4: image.Rect(300, 300, 600, 600),
		5: image.Rect(600, 300, 900, 600),
	}

	tileImages = map[int]*ebiten.Image{
		0: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[0]),
		1: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[1]),
		2: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[2]),
		3: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[3]),
		4: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[4]),
		5: getSubImageAndScaleDown(tilesheetImage, spriteMapRects[5]),
	}

	// println("Tile.init")
}

// Tile object describes a tile
type Tile struct {
	X, Y        int
	N, E, S, W  *Tile
	coins       uint
	tileImage   *ebiten.Image
	currDegrees int
	targDegrees int
	rotating    bool
	// if currDegrees != targDegrees then tile is lerping/rotating
	// rotating tile does not receive input

	// TODO map from coins bit field to sprite image index and rotate angle
	// TODO section this tile belongs to
	// TODO color
	// don't need hammingWeight because graphics not created dynamically
}

// NewTile creates a new Tile object and returns a pointer to it
func NewTile(x, y int) *Tile {
	t := &Tile{X: x, Y: y}
	return t
}

// PlaceCoin sets up a random config for this tile
func (t *Tile) PlaceCoin() {
	directions := [4]uint{NORTH, EAST, SOUTH, WEST}
	opposites := [4]uint{SOUTH, WEST, NORTH, EAST}
	links := [4]*Tile{t.N, t.E, t.S, t.W}
	// opplinks := [4]*Tile{t.S,t.W,t.N,t.E}
	for d := 0; d < 4; d++ {
		if rand.Float64() < 0.2 {
			if links[d] != nil {
				t.coins |= directions[d]
				links[d].coins |= opposites[d]
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
}

// Rect gives the x,y coords of the tile's top left and bottom right corners, in screen coordinates
func (t *Tile) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = t.X * TileWidth
	y0 = t.Y * TileHeight
	x1 = x0 + TileWidth
	y1 = y0 + TileHeight
	return // using named return parameters
}

// Rotate shifts the tile 90 degrees clockwise
func (t *Tile) Rotate() {
	if 0 == t.coins {
		return
	}
	if t.coins&WEST == WEST {
		// high bit is set
		t.coins = t.coins << 1
		t.coins = t.coins & MASK
		t.coins = t.coins | 1
	} else {
		t.coins = t.coins << 1
	}
	info := coin2ImageInfoMap[t.coins]
	t.targDegrees = info.deg
	// println("rotate", t.X, t.Y, t.coins)
}

// Update the tile state (transitions, user input)
func (t *Tile) Update() error {

	if t.currDegrees != t.targDegrees {
		t.currDegrees += 5
		if t.currDegrees >= 360 {
			t.currDegrees = 0
		}
	}

	return nil
}

// Draw handles rendering of Tile object
func (t *Tile) Draw(gridImage *ebiten.Image) error {

	// scale, point translation, rotate, object translation
	// extract sub image, scale it, rotate it, draw it into another image, then draw that constucted image into gridImage

	op := &ebiten.DrawImageOptions{}

	if t.currDegrees != 0 {
		op.GeoM.Translate(-float64(TileWidth)/2.0, -float64(TileHeight)/2.0)
		op.GeoM.Rotate(float64(float64(t.currDegrees) * 3.1415926 / float64(180)))
		op.GeoM.Translate(float64(TileWidth)/2.0, float64(TileHeight)/2.0)
	}

	x := t.X * TileWidth
	y := t.Y * TileHeight
	op.GeoM.Translate(float64(x), float64(y))
	gridImage.DrawImage(t.tileImage, op)

	return nil
}

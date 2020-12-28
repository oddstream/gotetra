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

var (
	spriteImages map[int]*ebiten.Image
	// TODO map from coins bit field to sprite image index and rotate angle
	// TODO state resting/rotating
	// TODO rotate lerp progress
	// TODO section this tile belongs to
	// TODO don't need hammingWeight because graphics not created dynamically
	// TODO color
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

	spriteImages = map[int]*ebiten.Image{
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
	X, Y       int
	N, E, S, W *Tile
	coins      uint
	tileImage  *ebiten.Image
}

// NewTile creates a new Tile object and returns a pointer to it
func NewTile(x, y int) *Tile {
	t := &Tile{X: x, Y: y}
	t.tileImage = spriteImages[rand.Intn(6)]
	return t
}

// Draw handles rendering of Tile object
func (t *Tile) Draw(gridImage *ebiten.Image) error {

	// scale, point translation, rotate, object translation
	// extract sub image, scale it, rotate it, draw it into another image, then draw that constucted image into gridImage

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-float64(TileWidth)/2.0, -float64(TileHeight)/2.0)
	op.GeoM.Rotate(float64(90.0 * 3.1415926 / 180))
	op.GeoM.Translate(float64(TileWidth)/2.0, float64(TileHeight)/2.0)

	x := t.X * TileWidth
	y := t.Y * TileHeight
	op.GeoM.Translate(float64(x), float64(y))
	gridImage.DrawImage(t.tileImage, op)

	return nil
}

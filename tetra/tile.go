// Copyright ©️ 2020-2021 oddstream.games

package tetra

import (
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// MinimumScale is the smallest a shape gets, used when creating and when completed
const MinimumScale float64 = 0.1

// NORTH is the bit pattern for the upwards direction
const (
	NORTH = 0b0001 // 1 << iota
	EAST  = 0b0010 // 1 << 1
	SOUTH = 0b0100 // 1 << 2
	WEST  = 0b1000 // 1 << 3
	MASK  = 0b1111
)

// TileState records what this tile is up to at the moment
type TileState int

// TileSettled is the state where a full sized tile is not doing anything
const (
	TileSettled TileState = iota
	TileGrowing
	TileShrinking
	TileRotating
	TileUnrotating
	TileShrunk
	TileBeingDragged
	TileReturning
)

var (
	tileImageLibrary map[uint]*ebiten.Image
	overSize         float64
	halfTileSize     float64
)

func initTileImages() {
	// used to be func init(), but TileSize may not be set yet, hence this func called from Grid init

	if 0 == TileSize {
		log.Fatal("Tile dimensions not set")
	}

	var makeFunc func(uint, int) image.Image = makeTileCurvy
	tileImageLibrary = make(map[uint]*ebiten.Image, 16)
	for i := uint(0); i < 16; i++ {
		img := makeFunc(i, TileSize)
		tileImageLibrary[i] = ebiten.NewImageFromImage(img)
	}

	// the tiles are all the same size, so pre-calc some useful variables
	actualTileSize, _ := tileImageLibrary[0].Size()
	halfTileSize = float64(actualTileSize) / 2
	overSize = float64((actualTileSize - TileSize) / 2)
}

// Tile object describes a tile
type Tile struct {
	G          *Grid
	X, Y       int
	N, E, S, W *Tile

	tileImage   *ebiten.Image
	currDegrees int
	targDegrees int // one of 0, 90, 180, 270
	scale       float64
	state       TileState
	coins       uint
	section     int
	color       color.RGBA

	homeX, homeY     float64 // position of tile when settled
	offsetX, offsetY float64 // offset to postion of tile when being dragged or lerping home

	// rotating, shrinking and growing tiles do not receive input
	// don't need hammingWeight
}

// NewTile creates a new Tile object and returns a pointer to it
// all new tiles start in a shrunken state
func NewTile(g *Grid, x, y int) *Tile {
	t := &Tile{G: g, X: x, Y: y, scale: MinimumScale, color: colorUnsectioned}
	// homeX, homeY, offsetX, offsetY will be (re)set by Layout()
	// coins and section will be 0
	return t
}

// Reset prepares a Tile for a new level by resetting just gameplay data, not structural data
func (t *Tile) Reset() {
	t.coins = 0
	t.section = 0
	t.color = colorUnsectioned //BasicColors["Black"]
	t.SetImage()               // reset to a blank tile image, will set state
}

// PlaceRandomCoins sets up a random config for this tile
func (t *Tile) PlaceRandomCoins() {
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
func (t *Tile) ColorConnected(colorName string, section int) {
	// println(colorName, ExtendedColors[colorName].R, ExtendedColors[colorName].G, ExtendedColors[colorName].B)
	bits := [4]uint{NORTH, EAST, SOUTH, WEST}
	links := [4]*Tile{t.N, t.E, t.S, t.W}

	t.color = ExtendedColors[colorName]
	t.section = section

	for d := 0; d < 4; d++ {
		if t.coins&bits[d] != 0 {
			tn := links[d]
			// unconnected tiles will have coins 0 and not have been colored (ie still be black)
			if tn != nil && tn.coins != 0 && tn.color == colorUnsectioned {
				tn.ColorConnected(colorName, section)
			}
		}
	}
}

// SetImage is used when all coins are placed
func (t *Tile) SetImage() {
	t.tileImage = tileImageLibrary[t.coins]
	if t.tileImage == nil {
		log.Fatal("tileImage is nil when coins == ", t.coins)
	}
	t.currDegrees = 0
	t.targDegrees = 0
	if t.coins == 0 {
		t.state = TileSettled
		t.scale = 1.0
	} else {
		t.state = TileGrowing
		t.scale = MinimumScale
	}
}

// Rect gives the x,y screen coords of the tile's top left and bottom right corners
func (t *Tile) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = t.X*TileSize + LeftMargin
	y0 = t.Y*TileSize + TopMargin
	x1 = x0 + TileSize
	y1 = y0 + TileSize
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
	if r < 0.1 {
		// rarely have to tap three times
		t.coins = shiftBits(t.coins)
	} else if r > 0.8 {
		// sometimes have to tap twice
		t.coins = shiftBits(t.coins)
		t.coins = shiftBits(t.coins)
	} else {
		// mostly have to tap once
		t.coins = unshiftBits(t.coins)
	}
	// TODO remove debugging jumble before release
	// if t.coins == NORTH || t.coins == SOUTH {
	// 	t.coins = unshiftBits(t.coins)
	// }
}

// if this tile had coins, how many points of connection would there be?
func (t Tile) bitsConnected(coins uint) int {
	var connected = func(dir uint, link *Tile, oppdir uint) int {
		if coins&dir == dir {
			if link != nil {
				if link.coins&oppdir == oppdir {
					return 1
				}
			}
		}
		return 0
	}
	var count int
	count += connected(NORTH, t.N, SOUTH)
	count += connected(EAST, t.E, WEST)
	count += connected(SOUTH, t.S, NORTH)
	count += connected(WEST, t.W, EAST)
	return count
}

// Rotate shifts the tile 90 degrees clockwise
func (t *Tile) Rotate() {
	if 0 == t.coins {
		return
	}
	if t.state != TileSettled {
		return
	}

	t.coins = shiftBits(t.coins)
	// switch t.currDegrees {
	// case 0:
	// 	t.targDegrees = 90
	// case 90:
	// 	t.targDegrees = 180
	// case 180:
	// 	t.targDegrees = 270
	// case 270:
	// 	t.targDegrees = 0
	// default:
	// 	println(t.currDegrees)
	// }
	t.targDegrees = (t.currDegrees + 90) % 360 // 360 % 360 == 0
	t.state = TileRotating
}

// Unrotate shifts the tile 90 degrees anticlockwise
func (t *Tile) Unrotate() {
	if 0 == t.coins {
		return
	}
	if t.state != TileSettled {
		return
	}

	t.coins = unshiftBits(t.coins)
	switch t.currDegrees {
	case 0:
		t.targDegrees = 270
	case 90:
		t.targDegrees = 0
	case 180:
		t.targDegrees = 90
	case 270:
		t.targDegrees = 180
	default:
		println(t.currDegrees)
	}
	// t.targDegrees = (t.currDegrees - 90) % 360 // -90 % 360 == 270 on calculator, but not in Go
	t.state = TileUnrotating
}

// IsComplete returns true if the tile aligns properly with it's neighbours
func (t *Tile) IsComplete() bool {
	if t.state != TileSettled {
		return false
	}
	if 0 == t.coins {
		return true
	}
	var test = func(dir uint, link *Tile, oppdir uint) bool {
		if (t.coins & dir) == dir {
			if (link == nil) || (link.coins&oppdir) == 0 {
				return false
			}
		}
		return true
	}
	if !test(NORTH, t.N, SOUTH) {
		return false
	}
	if !test(EAST, t.E, WEST) {
		return false
	}
	if !test(SOUTH, t.S, NORTH) {
		return false
	}
	if !test(WEST, t.W, EAST) {
		return false
	}
	return true
}

// IsSectionComplete returns true if the tile aligns properly with it's neighbours
func (t *Tile) IsSectionComplete(section int) bool {
	// By design, Go does not support optional parameters, default parameter values, or method overloading.
	if t.state != TileSettled {
		return false
	}
	if 0 == t.coins {
		return true
	}
	if section != t.section {
		return false
	}
	var test = func(dir uint, link *Tile, oppdir uint) bool {
		if (t.coins & dir) == dir {
			if (link == nil) || (link.section != section) || (link.coins&oppdir) == 0 {
				return false
			}
		}
		return true
	}
	if !test(NORTH, t.N, SOUTH) {
		return false
	}
	if !test(EAST, t.E, WEST) {
		return false
	}
	if !test(SOUTH, t.S, NORTH) {
		return false
	}
	if !test(WEST, t.W, EAST) {
		return false
	}
	return true
}

// TriggerScaleDown tells this tile to start scaling down
func (t *Tile) TriggerScaleDown() {
	if t.coins != 0 {
		t.state = TileShrinking
	}
}

// Layout the tile
func (t *Tile) Layout() {
	t.homeX = float64(LeftMargin + t.X*TileSize)
	t.homeY = float64(TopMargin + t.Y*TileSize)
	t.offsetX, t.offsetY = 0, 0
}

// Update the tile state (transitions, user input)
func (t *Tile) Update() error {

	if 0 == t.coins {
		return nil
	}

	switch t.state {
	case TileSettled:
		// nothing to do
	case TileGrowing:
		t.scale += 0.01
		if t.scale >= 1.0 {
			t.scale = 1.0
			t.state = TileSettled
		}
	case TileShrinking:
		t.scale -= 0.01
		if t.scale <= MinimumScale {
			t.scale = MinimumScale
			t.state = TileShrunk
		}
	case TileRotating:
		t.currDegrees = (t.currDegrees + 10) % 360
		if t.currDegrees == t.targDegrees {
			t.state = TileSettled
			if t.G.IsSectionComplete(t.section) {
				t.G.FilterSection((*Tile).TriggerScaleDown, t.section)
			}
		}
	case TileUnrotating:
		if t.currDegrees == 0 {
			t.currDegrees = 350
		} else {
			t.currDegrees -= 10
		}
		if t.currDegrees == t.targDegrees {
			t.state = TileSettled
			if t.G.IsSectionComplete(t.section) {
				t.G.FilterSection((*Tile).TriggerScaleDown, t.section)
			}
		}
	case TileShrunk:
		t.G.AddFrag(t.X, t.Y, t.tileImage, t.currDegrees, t.color)
		t.Reset()
	case TileBeingDragged:
		// nothing to do
	case TileReturning:
		// deliberately do a wrong lerp so tile starts fast and decelerates
		t.offsetX = lerp(t.offsetX, 0, 0.1)
		t.offsetY = lerp(t.offsetY, 0, 0.1)
		if t.offsetX < 1.0 && t.offsetY < 1.0 {
			t.offsetX, t.offsetY = 0, 0
			t.state = TileSettled
		}
	}

	return nil
}

func (t *Tile) debugText(screen *ebiten.Image, str string) {
	bound, _ := font.BoundString(Acme.small, str)
	w := (bound.Max.X - bound.Min.X).Ceil()
	h := (bound.Max.Y - bound.Min.Y).Ceil()
	x, y := t.homeX-overSize+t.offsetX, t.homeY-overSize+t.offsetY
	tx := int(x) + (TileSize-w)/2
	ty := int(y) + (TileSize-h)/2 + h
	var c color.Color = BasicColors["Black"]
	// if t.IsComplete() {
	// 	c = BasicColors["Fushia"]
	// } else {
	// 	c = BasicColors["Purple"]
	// }
	// ebitenutil.DrawRect(screen, float64(tx), float64(ty), float64(w), float64(h), c)
	text.Draw(screen, str, Acme.small, tx, ty, c)
}

// Draw handles rendering of Tile object
func (t *Tile) Draw(screen *ebiten.Image) {

	// scale, point translation, rotate, object translation

	op := &ebiten.DrawImageOptions{}

	if t.currDegrees != 0 {
		op.GeoM.Translate(-halfTileSize, -halfTileSize)
		op.GeoM.Rotate(float64(t.currDegrees) * 3.1415926535 / 180.0)
		op.GeoM.Translate(halfTileSize, halfTileSize)
	}

	// Reset RGB (not Alpha) forcibly
	// tilesheet already has black shapes
	if t.color != BasicColors["Black"] {
		// reducing alpha leaves the endcaps doubled
		op.ColorM.Scale(0, 0, 0, t.scale)
		// op.ColorM.Scale(0, 0, 0, 1)
		r := float64(t.color.R) / 0xff
		g := float64(t.color.G) / 0xff
		b := float64(t.color.B) / 0xff
		op.ColorM.Translate(r, g, b, 0)
		// op.CompositeMode = ebiten.CompositeModeLighter
	}

	if t.state == TileShrinking || t.state == TileShrunk || t.state == TileGrowing {
		// first move the origin to the center of the tile
		op.GeoM.Translate(-halfTileSize, -halfTileSize)
		op.GeoM.Scale(t.scale, t.scale)
		// then move the origin back to top left
		op.GeoM.Translate(halfTileSize, halfTileSize)
	}

	op.GeoM.Translate(t.homeX-overSize+t.offsetX, t.homeY-overSize+t.offsetY)

	// if t.X%2 == 0 && t.Y%2 == 0 {
	// 	colorTransBlack := color.RGBA{R: 0, G: 0, B: 0, A: 0x10}
	// 	ebitenutil.DrawRect(gridImage, float64(x), float64(y), float64(TileSize), float64(TileSize), colorTransBlack)
	// }

	screen.DrawImage(t.tileImage, op)

	// t.debugText(gridImage, fmt.Sprint(t.state))
	// t.debugText(gridImage, fmt.Sprintf("%04b", t.coins))
	// t.debugText(screen, fmt.Sprintf("%v,%v", t.currDegrees, t.targDegrees))
	// t.debugText(screen, fmt.Sprintf("%v", t.bitsConnected(t.coins)))
}

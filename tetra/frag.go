// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image/color"
	"math/rand"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

const fragSize float64 = 20.0
const halfFragSize float64 = fragSize / 2.0
const fragSizeInt int = 20

// Frag is an object that floats around the screen
type Frag struct {
	xCenter, yCenter float64
	dX, dY           float64       // direction
	rot              float64       // rotation
	rotVel           float64       // rotational velocity
	img              *ebiten.Image // image, scaled and colored
}

// NewFrag creates a new Frag and returns a pointer to it
func NewFrag(x, y float64, imgSrc *ebiten.Image, currDegrees float64, c color.RGBA) *Frag {
	f := &Frag{xCenter: x, yCenter: y, rot: currDegrees}

	f.rotVel = rand.Float64() - 0.5
	f.dX = rand.Float64() - 0.5
	f.dY = rand.Float64() - 0.5

	op := &ebiten.DrawImageOptions{}

	f.img = ebiten.NewImage(fragSizeInt, fragSizeInt)
	w, h := imgSrc.Size()
	scaleX := fragSize / float64(w)
	scaleY := fragSize / float64(h)
	op.GeoM.Scale(scaleX, scaleY)

	op.ColorM.Scale(0, 0, 0, 0.5)
	r := float64(c.R) / 0xff
	g := float64(c.G) / 0xff
	b := float64(c.B) / 0xff
	op.ColorM.Translate(r, g, b, 0)

	f.img.DrawImage(imgSrc, op)

	return f
}

// IsVisible returns true if frag is still visible
func (f *Frag) IsVisible() bool {
	var screenWidth, screenHeight int
	if runtime.GOARCH == "wasm" {
		screenWidth, screenHeight = WindowWidth, WindowHeight
	} else {
		screenWidth, screenHeight = ebiten.WindowSize()
	}
	return f.xCenter > 0 && f.xCenter < float64(screenWidth) && f.yCenter > 0 && f.yCenter < float64(screenHeight)
}

// Update the position of this Frag
func (f *Frag) Update() error {
	f.xCenter += f.dX
	f.yCenter += f.dY

	f.rot += f.rotVel
	if f.rot >= 360 {
		f.rot = 0
	}
	return nil
}

// Draw this Frag
func (f *Frag) Draw(gridImage *ebiten.Image) {

	sx, sy := f.img.Size()
	sxf, syf := float64(sx)/2.0, float64(sy)/2.0

	op := &ebiten.DrawImageOptions{}

	if f.rot != 0 {
		op.GeoM.Translate(-sxf, -syf)
		op.GeoM.Rotate(float64(f.rot) * 3.1415926535 / 180.0)
		op.GeoM.Translate(sxf, syf)
	}

	xOrigin := f.xCenter - sxf
	yOrigin := f.yCenter - syf

	op.GeoM.Translate(xOrigin, yOrigin)

	gridImage.DrawImage(f.img, op)
}

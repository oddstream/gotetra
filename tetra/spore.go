// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const sporeSize float64 = 20.0
const halfSporeSize float64 = sporeSize / 2.0
const sporeSizeInt int = 20

// Spore is an object that floats around the screen
type Spore struct {
	xCenter, yCenter float64
	dX, dY           float64       // direction
	rot              float64       // rotation
	rotVel           float64       // rotational velocy (-1, 0 or +1)
	img              *ebiten.Image // image, scaled and colored
}

// NewSpore creates a new Spore and returns a pointer to it
func NewSpore(x, y int, imgSrc *ebiten.Image, currDegrees float64, c color.RGBA) *Spore {
	sp := &Spore{xCenter: float64(x), yCenter: float64(y), rot: currDegrees}

	values := []float64{-1.0, 0.0, 1.0}
	sp.rotVel = values[rand.Intn(len(values))]
	sp.dX = (rand.Float64() - 0.5)
	sp.dY = (rand.Float64() - 0.5)

	op := &ebiten.DrawImageOptions{}

	sp.img = ebiten.NewImage(sporeSizeInt, sporeSizeInt)
	w, h := imgSrc.Size() // expecting this to be 100,100
	scaleX := sporeSize / float64(w)
	scaleY := sporeSize / float64(h)
	op.GeoM.Scale(scaleX, scaleY)

	op.ColorM.Scale(0, 0, 0, 0.5)
	r := float64(c.R) / 0xff
	g := float64(c.G) / 0xff
	b := float64(c.B) / 0xff
	op.ColorM.Translate(r, g, b, 0)

	sp.img.DrawImage(imgSrc, op)

	return sp
}

// IsVisible returns true of spore is still visible
func (sp *Spore) IsVisible() bool {
	return sp.xCenter > 0 && sp.xCenter < ScreenWidth && sp.yCenter > 0 && sp.yCenter < ScreenHeight
}

// Update the position of this Spore
func (sp *Spore) Update() error {
	sp.xCenter += sp.dX
	sp.yCenter += sp.dY

	sp.rot += sp.rotVel
	if sp.rot >= 360 {
		sp.rot = 0
	}
	return nil
}

// Draw this Spore
func (sp *Spore) Draw(gridImage *ebiten.Image) {

	sx, sy := sp.img.Size()
	sxf, syf := float64(sx)/2.0, float64(sy)/2.0

	op := &ebiten.DrawImageOptions{}

	if sp.rot != 0 {
		op.GeoM.Translate(-sxf, -syf)
		op.GeoM.Rotate(float64(sp.rot) * 3.1415926535 / 180.0)
		op.GeoM.Translate(sxf, syf)
	}

	xOrigin := sp.xCenter - sxf
	yOrigin := sp.yCenter - syf

	op.GeoM.Translate(xOrigin, yOrigin)

	gridImage.DrawImage(sp.img, op)
}

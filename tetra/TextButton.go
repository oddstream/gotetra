// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// TextButton is an object that represents a button
type TextButton struct {
	text             string
	font             font.Face
	xCenter, yCenter int
	xOrigin, yOrigin int
	width, height    int
}

// NewTextButton creates and returns a new TextButton object centered at x,y
func NewTextButton(str string, x, y int, btnFont font.Face) *TextButton {
	tb := &TextButton{text: str, xCenter: x, yCenter: y, font: btnFont}
	bound, _ := font.BoundString(tb.font, tb.text)
	tb.width = (bound.Max.X - bound.Min.X).Ceil()
	tb.height = (bound.Max.Y - bound.Min.Y).Ceil()
	tb.xOrigin = tb.xCenter - (tb.width / 2)
	tb.yOrigin = tb.yCenter - (tb.height / 2)
	return tb
}

// Pushed returns true if the button has just been pushed
func (tb *TextButton) Pushed() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > tb.xOrigin && x < tb.xOrigin+tb.width {
			if y > tb.yOrigin && y < tb.yOrigin+tb.height {
				return true
			}
		}
	}
	return false
}

// Update the button state (transitions, user input)
func (tb *TextButton) Update() error {
	return nil
}

// Draw handles rendering of TextButton object
func (tb *TextButton) Draw(screen *ebiten.Image) error {

	bgImage := ebiten.NewImage(tb.width, tb.height)
	bgImage.Fill(colorBlack)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(tb.xOrigin), float64(tb.yOrigin))
	screen.DrawImage(bgImage, op)

	text.Draw(screen, tb.text, tb.font, tb.xOrigin, tb.yOrigin+tb.height, colorWhite)

	return nil
}

// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// TextButton is an object that represents a button
type TextButton struct {
	text           string
	font           font.Face
	action         func()
	origin, center image.Point
	width, height  int
}

// NewTextButton creates and returns a new TextButton object centered at x,y
func NewTextButton(str string, x, y int, btnFont font.Face, actionFn func()) *TextButton {
	tb := &TextButton{text: str, center: image.Point{X: x, Y: y}, font: btnFont, action: actionFn}
	bound, _ := font.BoundString(tb.font, tb.text)
	tb.width = (bound.Max.X - bound.Min.X).Ceil()
	tb.height = (bound.Max.Y - bound.Min.Y).Ceil()
	tb.origin = image.Point{X: tb.center.X - (tb.width / 2), Y: tb.center.Y - (tb.height / 2)}
	return tb
}

// Rect gives the x,y coords of the TextButton's top left and bottom right corners, in screen coordinates
func (tb *TextButton) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = tb.origin.X
	y0 = tb.origin.Y
	x1 = x0 + tb.width
	y1 = y0 + tb.height
	return // using named return parameters
}

// Pushed returns true if the button has just been pushed
func (tb *TextButton) Pushed(i *Input) bool {
	if i.X != 0 && i.Y != 0 {
		return InRect(i.X, i.Y, tb.Rect)
	}
	return false
}

// Action invikes the action func
func (tb *TextButton) Action() {
	if tb.action != nil {
		tb.action()
	}
}

// Update the button state (transitions, NOT user input)
func (tb *TextButton) Update() error {
	return nil
}

// Draw handles rendering of TextButton object
func (tb *TextButton) Draw(screen *ebiten.Image) {

	bgImage := ebiten.NewImage(tb.width, tb.height)
	bgImage.Fill(BasicColors["Black"])
	op := &ebiten.DrawImageOptions{}
	{
		op.GeoM.Translate(-float64(tb.width)/2.0, -float64(tb.height)/2.0)
		op.GeoM.Scale(1.1, 1.5)
		op.GeoM.Translate(float64(tb.width)/2.0, float64(tb.height)/2.0)
	}
	op.GeoM.Translate(float64(tb.origin.X), float64(tb.origin.Y))
	screen.DrawImage(bgImage, op)

	text.Draw(screen, tb.text, tb.font, tb.origin.X, tb.origin.Y+tb.height, BasicColors["White"])

}

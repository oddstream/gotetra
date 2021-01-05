// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Label is an object that represents a button
type Label struct {
	text           string
	font           font.Face
	origin, center image.Point
	width, height  int
}

// NewLabel creates and returns a new Label object centered at x,y
func NewLabel(str string, x, y int, btnFont font.Face) *Label {
	l := &Label{text: str, center: image.Point{X: x, Y: y}, font: btnFont}
	bound, _ := font.BoundString(l.font, l.text)
	l.width = (bound.Max.X - bound.Min.X).Ceil()
	l.height = (bound.Max.Y - bound.Min.Y).Ceil()
	l.origin = image.Point{X: l.center.X - (l.width / 2), Y: l.center.Y - (l.height / 2)}
	return l
}

// Rect gives the x,y coords of the label's top left and bottom right corners, in screen coordinates
func (l *Label) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = l.origin.X
	y0 = l.origin.Y
	x1 = x0 + l.width
	y1 = y0 + l.height
	return // using named return parameters
}

// Pushed returns true if the label has just been pushed, which we don't care about
func (l *Label) Pushed(*Input) bool {
	return false
}

// Action invikes the action func
func (l *Label) Action() {
	// Labels take no action
}

// Update the button state (transitions, NOT user input)
func (l *Label) Update() error {
	return nil
}

// Draw handles rendering of Label object
func (l *Label) Draw(screen *ebiten.Image) {

	text.Draw(screen, l.text, l.font, l.origin.X, l.origin.Y+l.height, BasicColors["White"])

}

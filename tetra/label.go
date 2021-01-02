// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Label is an object that represents a button
type Label struct {
	text             string
	font             font.Face
	xCenter, yCenter int
	xOrigin, yOrigin int
	width, height    int
}

// NewLabel creates and returns a new Label object centered at x,y
func NewLabel(str string, x, y int, btnFont font.Face) *Label {
	tb := &Label{text: str, xCenter: x, yCenter: y, font: btnFont}
	bound, _ := font.BoundString(tb.font, tb.text)
	tb.width = (bound.Max.X - bound.Min.X).Ceil()
	tb.height = (bound.Max.Y - bound.Min.Y).Ceil()
	tb.xOrigin = tb.xCenter - (tb.width / 2)
	tb.yOrigin = tb.yCenter - (tb.height / 2)
	return tb
}

// Pushed returns true if the label has just been pushed, which we don't care about
func (tb *Label) Pushed() bool {
	return false
}

// Action invikes the action func
func (tb *Label) Action() {
	// Labels take no action
}

// Update the button state (transitions, user input)
func (tb *Label) Update() error {
	return nil
}

// Draw handles rendering of Label object
func (tb *Label) Draw(screen *ebiten.Image) {

	text.Draw(screen, tb.text, tb.font, tb.xOrigin, tb.yOrigin+tb.height, BasicColors["White"])

}

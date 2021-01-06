// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Input records state of mouse and touch
type Input struct {
	pt image.Point
}

// NewInput Input object constructor
func NewInput() *Input {
	// no fields to initialize, so use the built-in new()
	return new(Input)
}

// Update the state of the Input object
func (i *Input) Update() {
	x, y := 0, 0
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y = ebiten.CursorPosition()
	}
	ts := inpututil.JustPressedTouchIDs()
	if ts != nil && len(ts) == 1 {
		if inpututil.IsTouchJustReleased(ts[0]) {
			x, y = ebiten.TouchPosition(ts[0])
		}
	}
	i.pt = image.Point{X: x, Y: y}
}

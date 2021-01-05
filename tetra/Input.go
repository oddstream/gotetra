// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Input records state of mouse and touch
type Input struct {
	X, Y int
}

// NewInput Input object constructor
func NewInput() *Input {
	// no fields to initialize, so use the built-in new()
	return new(Input)
}

// Update the state of the Input object
func (i *Input) Update() {
	i.X, i.Y = 0, 0
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		i.X, i.Y = ebiten.CursorPosition()
	}
	ts := inpututil.JustPressedTouchIDs()
	if ts != nil && len(ts) == 1 {
		if inpututil.IsTouchJustReleased(ts[0]) {
			i.X, i.Y = ebiten.TouchPosition(ts[0])
		}
	}
}

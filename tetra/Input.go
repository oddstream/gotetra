// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// var (
// 	Nothing = struct{}{}
// )

// StrokeSource represents a input device to provide strokes.
type StrokeSource interface {
	Position() (int, int)
	IsJustReleased() bool
}

// MouseStrokeSource is a StrokeSource implementation of mouse.
type MouseStrokeSource struct{}

// Position returns the x,y cordinates of the cursor position
func (m *MouseStrokeSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

// IsJustReleased returns true if the left mouse button was released in the current frame
func (m *MouseStrokeSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

// TouchStrokeSource is a StrokeSource implementation of touch.
type TouchStrokeSource struct {
	ID ebiten.TouchID
}

// Position returns the x,y cordinates of the cursor position
func (t *TouchStrokeSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}

// IsJustReleased returns true if the first touch was released in the current frame
func (t *TouchStrokeSource) IsJustReleased() bool {
	return inpututil.IsTouchJustReleased(t.ID)
}

// Stroke manages the current drag state by mouse.
type Stroke struct {
	source StrokeSource

	// initX and initY represents the position when dragging starts.
	initX, initY int

	// currentX and currentY represents the current position
	currentX, currentY int

	released bool

	// draggingObject represents a object (like a tile) that is being dragged.
	draggingObject interface{}
}

// NewStroke creates a new Stroke object
func NewStroke(source StrokeSource) *Stroke {
	cx, cy := source.Position()
	return &Stroke{
		source:   source,
		initX:    cx,
		initY:    cy,
		currentX: cx,
		currentY: cy,
	}
}

// Update is called once per frame and updates the Stroke object
func (s *Stroke) Update() {
	if s.released {
		return
	}
	if s.source.IsJustReleased() {
		s.released = true
		return
	}
	x, y := s.source.Position()
	s.currentX = x
	s.currentY = y
}

// IsReleased returns true if ...
func (s *Stroke) IsReleased() bool {
	return s.released
}

// Position returns the x,y position of the cursor
func (s *Stroke) Position() (int, int) {
	return s.currentX, s.currentY
}

// PositionPoint returns the x,y position of the cursor as a point
func (s *Stroke) PositionPoint() image.Point {
	pt := image.Point{X: s.currentX, Y: s.currentY}
	return pt
}

// PositionDiff returns the x,y difference between the start of the stroke and the stoke's current position
func (s *Stroke) PositionDiff() (int, int) {
	dx := s.currentX - s.initX
	dy := s.currentY - s.initY
	return dx, dy
}

// DraggingObject returns a reference to the object currently being dragged
func (s *Stroke) DraggingObject() interface{} {
	return s.draggingObject
}

// SetDraggingObject sets the object currently being dragged
func (s *Stroke) SetDraggingObject(object interface{}) {
	s.draggingObject = object
}

// Input records state of mouse and touch
type Input struct {
	pt          image.Point
	backPressed bool
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

	i.backPressed = inpututil.IsKeyJustReleased(ebiten.KeyBackspace)
}

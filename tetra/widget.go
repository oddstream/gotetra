package tetra

import "github.com/hajimehoshi/ebiten/v2"

// Widget type implements UpDate, Draw and Pushed
type Widget interface {
	Update() error
	Draw(*ebiten.Image)
	SetPosition(int, int)
	Rect() (int, int, int, int)
	Pushed(*Input) bool
	Action()
}

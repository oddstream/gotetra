// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Drawable type implements UpDate, Draw and Pushed
type Drawable interface {
	Update() error
	Draw(*ebiten.Image)
	Pushed() bool
	Action()
}

// Splash represents a game state.
type Splash struct {
	logoImage  *ebiten.Image
	xPos, yPos int
	widgets    []Drawable
}

// NewSplash creates and initializes a Splash/GameState object
func NewSplash() *Splash {
	s := &Splash{}
	var err error
	s.logoImage, _, err = ebitenutil.NewImageFromFile("/home/gilbert/Tetra/assets/oddstream logo.png")
	if err != nil {
		log.Fatal(err)
	}
	sx, sy := s.logoImage.Size()
	s.xPos = (ScreenWidth - sx) / 2
	s.yPos = -sy

	s.widgets = []Drawable{
		NewLabel("Do you prefer", ScreenWidth/2, 400, Acme.normal),
		NewTextButton("LITTLE PUZZLES", ScreenWidth/2, 500, Acme.large, func() { GSM.Switch(NewPuzzle("puzzle", 4, 5)) }),
		NewLabel("or", ScreenWidth/2, 600, Acme.normal),
		NewTextButton("BUBBLE WRAP", ScreenWidth/2, 700, Acme.large, func() { GSM.Switch(NewPuzzle("bubblewrap", 4, 5)) }),
	}
	return s
}

// Layout implements ebiten.Game's Layout
func (s *Splash) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (s *Splash) Update() error {
	if s.yPos < 50 {
		s.yPos += ScreenWidth / ebiten.DefaultTPS
	}

	for _, w := range s.widgets {
		if w.Pushed() {
			w.Action()
			break
		}
	}

	return nil
}

// Draw draws the current GameState to the given screen
func (s *Splash) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.xPos), float64(s.yPos))
	screen.DrawImage(s.logoImage, op)

	for _, d := range s.widgets {
		d.Draw(screen)
	}

	ebitenutil.DrawLine(screen, 0, 500, ScreenWidth, 500, BasicColors["Black"])
	ebitenutil.DrawLine(screen, 0, 700, ScreenWidth, 700, BasicColors["Black"])
}

// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"bytes"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Widget type implements UpDate, Draw and Pushed
type Widget interface {
	Update() error
	Draw(*ebiten.Image)
	Rect() (int, int, int, int)
	Pushed(*Input) bool
	Action()
}

// Pushable type implements Rect
type Pushable interface {
	Rect() (int, int, int, int)
}

// Splash represents a game state.
type Splash struct {
	logoImage *ebiten.Image
	pos       image.Point
	widgets   []Widget
	input     *Input
}

// NewSplash creates and initializes a Splash/GameState object
func NewSplash() *Splash {
	s := &Splash{input: NewInput()}
	// var err error
	// s.logoImage, _, err = ebitenutil.NewImageFromFile("/home/gilbert/Tetra/assets/oddstream logo.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Decode image from a byte slice instead of a file so that this works in any working directory.
	// If you want to use a file, there are some options:
	// 1) Use os.Open and pass the file to the image decoder.
	//    This is a very regular way, but doesn't work on browsers.
	// 2) Use ebitenutil.OpenFile and pass the file to the image decoder.
	//    This works even on browsers.
	// 3) Use ebitenutil.NewImageFromFile to create an ebiten.Image directly from a file.
	//    This also works on browsers.
	img, _, err := image.Decode(bytes.NewReader(Logo_png))
	if err != nil {
		log.Fatal(err)
	}
	s.logoImage = ebiten.NewImageFromImage(img)

	sx, sy := s.logoImage.Size()
	s.pos = image.Point{X: (ScreenWidth - sx) / 2, Y: -sy}

	xCenter := ScreenWidth / 2

	yPlaces := [6]int{}
	for i := 0; i < 6; i++ {
		yPlaces[i] = (ScreenHeight / 6) * i
	}
	s.widgets = []Widget{
		NewTextButton("EASY", xCenter, yPlaces[2], Acme.large, func() { GSM.Switch(NewGrid("bubblewrap", 7, 6)) }),
		NewTextButton("NORMAL", xCenter, yPlaces[3], Acme.large, func() { GSM.Switch(NewGrid("bubblewrap", 0, 0)) }),
		NewTextButton("HARD", xCenter, yPlaces[4], Acme.large, func() { GSM.Switch(NewGrid("puzzle", 0, 0)) }),
		NewTextButton("HARDEST", xCenter, yPlaces[5], Acme.large, func() { GSM.Switch(NewGrid("puzzle", 18, 10)) }),
	}
	return s
}

// Layout implements ebiten.Game's Layout
func (s *Splash) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (s *Splash) Update() error {

	if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) {
		os.Exit(0)
	}

	s.input.Update()

	if s.pos.Y < 0 {
		s.pos.Y++
	}

	for _, w := range s.widgets {
		if w.Pushed(s.input) {
			w.Action()
			break
		}
	}

	return nil
}

// Draw draws the current GameState to the given screen
func (s *Splash) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.pos.X), float64(s.pos.Y))
	screen.DrawImage(s.logoImage, op)

	for _, d := range s.widgets {
		d.Draw(screen)
	}

	// ebitenutil.DrawLine(screen, 0, 500, ScreenWidth, 500, BasicColors["Black"])
	// ebitenutil.DrawLine(screen, 0, 700, ScreenWidth, 700, BasicColors["Black"])
}

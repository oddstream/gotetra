// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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
	s.xPos = (ScreenWidth - sx) / 2
	s.yPos = -sy

	xCenter := ScreenWidth / 2

	s.widgets = []Drawable{
		NewLabel("Do you prefer", xCenter, 200, Acme.normal),
		NewTextButton("LITTLE PUZZLES", xCenter, 300, Acme.large, func() { GSM.Switch(NewGrid("puzzle")) }),
		NewLabel("or", xCenter, 400, Acme.normal),
		NewTextButton("BUBBLE WRAP", xCenter, 500, Acme.large, func() { GSM.Switch(NewGrid("bubblewrap")) }),
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

	// ebitenutil.DrawLine(screen, 0, 500, ScreenWidth, 500, BasicColors["Black"])
	// ebitenutil.DrawLine(screen, 0, 700, ScreenWidth, 700, BasicColors["Black"])
}

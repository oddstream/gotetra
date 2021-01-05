// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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
	logoImage  *ebiten.Image
	xPos, yPos int
	widgets    []Widget
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

	yPlaces := [6]int{}
	for i := 0; i < 6; i++ {
		yPlaces[i] = (ScreenHeight / 6) * i
	}
	s.widgets = []Widget{
		NewLabel("Do you prefer", xCenter, yPlaces[2], Acme.normal),
		NewTextButton("LITTLE PUZZLES", xCenter, yPlaces[3], Acme.large, func() { GSM.Switch(NewGrid("puzzle", 7, 7)) }),
		NewLabel("or", xCenter, yPlaces[4], Acme.normal),
		NewTextButton("BUBBLE WRAP", xCenter, yPlaces[5], Acme.large, func() { GSM.Switch(NewGrid("bubblewrap", 0, 0)) }),
	}
	return s
}

// Layout implements ebiten.Game's Layout
func (s *Splash) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (s *Splash) Update(i *Input) error {

	i.Update()

	if s.yPos < 0 {
		s.yPos++
	}

	for _, w := range s.widgets {
		if w.Pushed(i) {
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
	op.GeoM.Translate(float64(s.xPos), float64(s.yPos))
	screen.DrawImage(s.logoImage, op)

	for _, d := range s.widgets {
		d.Draw(screen)
	}

	// ebitenutil.DrawLine(screen, 0, 500, ScreenWidth, 500, BasicColors["Black"])
	// ebitenutil.DrawLine(screen, 0, 700, ScreenWidth, 700, BasicColors["Black"])
}

package tetra

import (
	"bytes"
	_ "embed" // go:embed only allowed in Go files that import "embed"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/raccoon280x180.png
var logoBytes []byte

// Splash represents a game state.
type Splash struct {
	logoImage *ebiten.Image
	logoPos   image.Point
	widgets   []Widget
	input     *Input
}

// NewSplash creates and initializes a Splash/GameState object
func NewSplash() *Splash {
	s := &Splash{input: NewInput()}
	img, _, err := image.Decode(bytes.NewReader(logoBytes))
	if err != nil {
		log.Fatal(err)
	}
	s.logoImage = ebiten.NewImageFromImage(img)

	s.widgets = []Widget{
		NewLabel("T E T R A                      L O O P S", Acme.large),
		NewTextButton("EASY", Acme.large, func() { GSM.Switch(NewGrid("bubblewrap", 6, 5)) }),
		NewTextButton("NORMAL", Acme.large, func() { GSM.Switch(NewGrid("bubblewrap", 0, 0)) }),
		NewTextButton("HARD", Acme.large, func() { GSM.Switch(NewGrid("puzzle", 0, 0)) }),
		NewTextButton("HARDEST", Acme.large, func() { GSM.Switch(NewGrid("puzzle", 14, 7)) }),
	}

	return s
}

// Layout implements ebiten.Game's Layout
func (s *Splash) Layout(outsideWidth, outsideHeight int) (int, int) {

	xCenter := outsideWidth / 2
	// create 6 vertical slots for 5 widgets
	yPlaces := [6]int{} // golang gotcha: can't use len(s.widgets)
	for i := 0; i < len(yPlaces); i++ {
		yPlaces[i] = (outsideHeight / len(yPlaces)) * i
	}

	lx, ly := s.logoImage.Size()
	s.logoPos = image.Point{X: xCenter - (lx / 2), Y: yPlaces[1] - (ly / 2)}

	for i, w := range s.widgets {
		w.SetPosition(xCenter, yPlaces[i+1])
	}

	return outsideWidth, outsideHeight
}

// Update updates the current game state.
func (s *Splash) Update() error {

	if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) {
		os.Exit(0)
	}

	s.input.Update()

	for i, w := range s.widgets {
		if w.Pushed(s.input) {
			PlayPianoNote(i)
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
	{
		sx, sy := s.logoImage.Size()
		sx, sy = sx/2, sy/2
		op.GeoM.Translate(float64(-sx), float64(-sy))
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(float64(sx), float64(sy))

		op.GeoM.Translate(float64(s.logoPos.X), float64(s.logoPos.Y))
		screen.DrawImage(s.logoImage, op)
	}

	for _, d := range s.widgets {
		d.Draw(screen)
	}

	// ebitenutil.DrawLine(screen, 0, 500, ScreenWidth, 500, BasicColors["Black"])
	// ebitenutil.DrawLine(screen, 0, 700, ScreenWidth, 700, BasicColors["Black"])
}

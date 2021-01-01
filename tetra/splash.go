// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Splash represents a game state.
type Splash struct {
	logoImage  *ebiten.Image
	xPos, yPos int
	btnPuzzle  *TextButton
	btnBubble  *TextButton
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

	s.btnPuzzle = NewTextButton("LITTLE PUZZLES", ScreenWidth/2, 500, Acme.large)
	s.btnBubble = NewTextButton("BUBBLE WRAP", ScreenWidth/2, 700, Acme.large)
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
	/*
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			_, y := ebiten.CursorPosition()
			if y < (ScreenHeight / 2) {
				GSM.Switch(NewPuzzle("puzzle", 4, 7))
			} else {
				GSM.Switch(NewPuzzle("bubblewrap", 6, 9))
			}
		}
	*/
	if s.btnPuzzle.Pushed() {
		GSM.Switch(NewPuzzle("puzzle", 4, 7))
	} else if s.btnBubble.Pushed() {
		GSM.Switch(NewPuzzle("bubblewrap", 6, 9))
	}
	return nil
}

// Draw draws the current GameState to the given screen
func (s *Splash) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.xPos), float64(s.yPos))
	screen.DrawImage(s.logoImage, op)

	s.btnPuzzle.Draw(screen)
	s.btnBubble.Draw(screen)

	ebitenutil.DrawLine(screen, 0, 500, ScreenWidth, 500, colorBlack)
	ebitenutil.DrawLine(screen, 0, 700, ScreenWidth, 700, colorBlack)

}

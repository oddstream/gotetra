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
}

// Init initializes a Splash/GameState object that was created by the caller
func (s *Splash) Init() {
	var err error
	s.logoImage, _, err = ebitenutil.NewImageFromFile("assets/oddstream logo.png")
	if err != nil {
		log.Fatal(err)
	}
	sx, sy := s.logoImage.Size()
	s.xPos = -sx
	s.yPos = (ScreenHeight - sy) / 2
}

// Layout implements ebiten.Game's Layout
func (s *Splash) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (s *Splash) Update() error {
	s.xPos += 20
	if s.xPos > ScreenWidth {
		println("change state to puzzle")
		pz := &Puzzle{}
		pz.Init()
		GSM.Switch(pz)
	}
	return nil
}

// Draw draws the current GameState to the given screen
func (s *Splash) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.xPos), float64(s.yPos))
	screen.DrawImage(s.logoImage, op)
}

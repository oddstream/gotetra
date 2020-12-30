// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ScreenWidth and ScreenHeight are exported to main.go for ebiten.SetWindowSize()
const (
	ScreenWidth  = 700
	ScreenHeight = 1000
)

// Game represents a game state.
type Game struct {
}

// GSM provides global access to the game state manager
var GSM *GameStateManager = &GameStateManager{}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{}

	GSM.Switch(NewSplash())

	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	state := GSM.Get()
	return state.Layout(ScreenWidth, ScreenHeight)
}

// Update updates the current game state.
func (g *Game) Update() error {
	state := GSM.Get()
	if err := state.Update(); err != nil {
		return err
	}
	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	state := GSM.Get()
	state.Draw(screen)
}

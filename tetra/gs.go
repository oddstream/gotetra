// Copyright ©️ 2020 oddstream.games

package tetra

import "github.com/hajimehoshi/ebiten/v2"

// GameState interface defines the API for each game state
// each seperate game state (eg Splash, Menu, Puzzle, GameOver) must implement these
type GameState interface {
	Layout(int, int) (int, int)
	Update() error
	Draw(*ebiten.Image)
}

// Copyright ©️ 2020 oddstream.games

package tetra

// GameStateManager does what it says on the tin
type GameStateManager struct {
	currentState GameState
}

// Switch changes to a different GameState
func (gsm *GameStateManager) Switch(state GameState) {
	gsm.currentState = state
}

// Get returns the current GameState
func (gsm *GameStateManager) Get() GameState {
	return gsm.currentState
}

// Copyright ©️ 2020 oddstream.games

// go mod init oddstream.games/tetra
// go mod tidy

// the package defining a command (an excutable Go program) always has the name main
// this is a signal to go build that it must invoke the linker to make an executable file
package main

import (
	"log"

	// load png decoder in main package
	// _ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	tetra "oddstream.games/tetra/tetra"
)

func main() {
	game, err := tetra.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowTitle("Tetra Loops")
	ebiten.SetWindowSize(1920/2, 1080/2)
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

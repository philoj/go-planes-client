package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten"
	"log"
	"watchYourSix/lobby"
	"watchYourSix/planes"
)

var playerId = flag.Int("id", 1, "Set a unique id for each client")

func main() {
	flag.Parse()
	log.Println("Player id:", *playerId)
	game := planes.NewGame(*playerId)
	// server
	go lobby.JoinLobby(game)

	// this is for testing only
	ebiten.SetRunnableOnUnfocused(true)

	//the game ui

	ebiten.SetWindowSize(planes.ScreenWidth, planes.ScreenHeight)
	ebiten.SetWindowTitle("Watch you Six")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

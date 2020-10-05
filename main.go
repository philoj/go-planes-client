package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten"
	"log"
	"watchYourSix/planes"
)

var playerId = flag.Int("id", 1, "Set a unique id for each client")
var screenWidth = flag.Int("w", 600, "Screen Width in pixels")
var screenHeight = flag.Int("h", 600, "Screen height in pixels")
var debug = flag.Bool("debug", false, "Debug enabled(default false)")
func main() {
	flag.Parse()
	log.Println("Player id:", *playerId)
	game := planes.NewGame(*debug)

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(*screenWidth, *screenHeight)
	ebiten.SetWindowTitle("Watch you Six")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

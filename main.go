package main

import (
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"goplanesclient/game"
	"log"
)

var playerId = flag.Int("id", 1, "Set a unique id for each client") // FIXME: generate unique id if omitted
var screenWidth = flag.Int("w", 600, "Screen Width in pixels")
var screenHeight = flag.Int("h", 600, "Screen height in pixels")
var debug = flag.Bool("debug", false, "Debug enabled(default false)")
var host = flag.String("host", "localhost:8080", "Debug enabled(default false)")
var path = flag.String("path", "/lobby", "Debug enabled(default false)")

func main() {
	flag.Parse()
	log.Println("Player id:", *playerId)
	*path = fmt.Sprintf("%s/%d", *path, *playerId)
	g := game.NewGame(*playerId, *debug, *host, *path)
	configureEbiten()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func configureEbiten() {
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(*screenWidth, *screenHeight)
	ebiten.SetWindowTitle("Watch you Six")
}

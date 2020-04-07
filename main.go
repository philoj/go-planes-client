package main

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"watchYourSix/planes"
)


func main() {
	ebiten.SetWindowSize(planes.ScreenWidth*2, planes.ScreenHeight*2)
	ebiten.SetWindowTitle("Watch you Six")
	if err := ebiten.RunGame(planes.NewGame()); err != nil {
		log.Fatal(err)
	}
}

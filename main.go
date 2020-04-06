package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten"
	"log"
	"os"
	"runtime/pprof"
	"watchYourSix/planes"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}
	ebiten.SetWindowSize(planes.ScreenWidth*2, planes.ScreenHeight*2)
	ebiten.SetWindowTitle("Watch you Six")
	if err := ebiten.RunGame(planes.NewGame()); err != nil {
		log.Fatal(err)
	}
}

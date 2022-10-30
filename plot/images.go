package plot

import (
	"github.com/hajimehoshi/ebiten/v2"
	"goplanesclient/geometry"
	"math"
)

func DrawImage(screen, img *ebiten.Image, translate geometry.Vector, heading float64) {
	w, h := img.Size()
	rotScale := &ebiten.DrawImageOptions{}
	rotScale.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	rotScale.GeoM.Rotate(heading)
	rotScale.GeoM.Translate(translate.I, translate.J)
	screen.DrawImage(img, rotScale)
}
func LaySquareTiledImage(screen, tile *ebiten.Image, originalTranslation geometry.Vector, tileSize float64, offsetCount int) {
	bgOptions := &ebiten.DrawImageOptions{}
	oc := float64(offsetCount)
	equivalentTranslation := geometry.Vector{
		I: math.Mod(originalTranslation.I, tileSize),
		J: math.Mod(originalTranslation.J, tileSize),
	}.Add(geometry.Vector{
		I: tileSize * oc,
		J: tileSize * oc,
	})
	bgOptions.GeoM.Translate(equivalentTranslation.I, equivalentTranslation.J)
	screen.DrawImage(tile, bgOptions)
}

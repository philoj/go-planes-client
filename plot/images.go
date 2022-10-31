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

func LaySquareTiles(screen, tile *ebiten.Image, originalTranslation geometry.Vector) {
	w, h := screen.Size()
	tileSize := float64(tile.Bounds().Dx())
	dx, dy := math.Mod(originalTranslation.I, tileSize), math.Mod(originalTranslation.J, tileSize)
	for x := -tileSize; x <= float64(w)+tileSize; x += tileSize {
		for y := -tileSize; y <= float64(h)+tileSize; y += tileSize {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(x+dx, y+dy)
			screen.DrawImage(tile, opt)
		}
	}
}

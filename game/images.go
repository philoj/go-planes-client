package game

import "github.com/hajimehoshi/ebiten"

type imageSize struct {
	width, height int
}
type imageInfo struct {
	path                     string
	originalSize, targetSize imageSize
	image                    *ebiten.Image
}

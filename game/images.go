package game

import "github.com/hajimehoshi/ebiten/v2"

type imageSize struct {
	width, height int
}
type imageInfo struct {
	path                     string
	originalSize, targetSize imageSize
	image                    *ebiten.Image
}

var (
	images = map[string]*imageInfo{
		bgImageAssetId: {
			path: "/bg.jpg",
			originalSize: imageSize{
				width:  bgImageSize,
				height: bgImageSize,
			},
		},
		iconImageAssetId: {
			path: "/icon_orig.png",
			originalSize: imageSize{
				width:  playerIconImageSize,
				height: playerIconImageSize,
			},
		},
		blipImageAssetId: {
			path: "/blip.png",
			originalSize: imageSize{
				width:  blipIconImageSize,
				height: blipIconImageSize,
			},
		},
	}
)

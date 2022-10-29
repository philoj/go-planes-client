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

var (
	images = map[string]*imageInfo{
		BgImageAssetId: {
			path: "/bg.jpg",
			originalSize: imageSize{
				width:  bgImageSize,
				height: bgImageSize,
			},
		},
		IconImageAssetId: {
			path: "/icon_orig.png",
			originalSize: imageSize{
				width:  playerIconImageSize,
				height: playerIconImageSize,
			},
		},
		BlipImageAssetId: {
			path: "/blip.png",
			originalSize: imageSize{
				width:  blipIconImageSize,
				height: blipIconImageSize,
			},
		},
	}
)

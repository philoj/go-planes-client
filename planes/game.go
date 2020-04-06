package planes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	_ "image/jpeg"
	"math"
)

const (
	ScreenWidth        = 300.0
	ScreenHeight       = 300.0
	bgImageHeight      = 5000.0
	bgImageWidth       = 5000.0
	bgScaleX           = ScreenWidth / bgImageWidth
	bgScaleY           = ScreenHeight / bgImageHeight
	seamCountY         = float64(2)
	seamCountX         = float64(2)
	seamOffsetLimitX   = ScreenWidth
	seamOffsetLimitY   = ScreenHeight
	initialSeamOffsetY = -ScreenHeight
	initialSeamOffsetX = -ScreenWidth

	playerIconScaleX = 0.05
	playerIconScaleY = 0.05

	initialVelocity = 2.0
)

var bgLayer, playerIcon *ebiten.Image

func init() {
	tile, _, _ := ebitenutil.NewImageFromFile("planes/assets/bg.jpg", ebiten.FilterDefault)
	player, _, _ := ebitenutil.NewImageFromFile("planes/assets/icon.jpg", ebiten.FilterDefault)
	playerIcon, _ = ebiten.NewImage(32, 32, ebiten.FilterDefault)
	pTran := ebiten.ScaleGeo(playerIconScaleX, playerIconScaleY)
	playerIcon.DrawImage(player, &ebiten.DrawImageOptions{
		GeoM:          pTran,
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
		ImageParts:    nil,
		Parts:         nil,
		SourceRect:    nil,
	})

	tileCountX := int(seamCountX + 1)
	tileCountY := int(seamCountY + 1)
	bgLayer, _ = ebiten.NewImage(ScreenWidth*tileCountX, ScreenHeight*tileCountY, ebiten.FilterDefault)

	for i := -0; i < tileCountY; i++ {
		for j := -0; j < tileCountX; j++ {
			bgTileTransform := ebiten.ScaleGeo(bgScaleX, bgScaleY)
			bgTileTransform.Translate(float64(j)*ScreenWidth, float64(i)*ScreenHeight)
			fmt.Println(tile)
			bgLayer.DrawImage(tile, &ebiten.DrawImageOptions{
				GeoM:          bgTileTransform,
				ColorM:        ebiten.ColorM{},
				CompositeMode: 0,
				Filter:        0,
				ImageParts:    nil,
				Parts:         nil,
				SourceRect:    nil,
			})
		}
	}
}

type Game struct {
	x              float64
	y              float64
	bgLayerOffsetY float64
}

func NewGame() *Game {
	return &Game{
		bgLayerOffsetY: 0.0,
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	dx, dy := 0.0, initialVelocity
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && ebiten.IsKeyPressed(ebiten.KeyRight) {
		fmt.Println("both")
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		fmt.Println("left")
		dx -= 0.1 * initialVelocity
		dy -= initialVelocity/4
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		fmt.Println("Right")
		dx += 0.1 * initialVelocity
		dy -= initialVelocity/4
	}
	g.move(screen, dx, dy)
	return nil
}

func (g *Game) move(screen *ebiten.Image, dx float64, dy float64) {
	g.x += dx
	g.y += dy
	bgOffsetX := math.Mod(g.x, seamOffsetLimitX)
	bgOffsetY := math.Mod(g.y, seamOffsetLimitY)
	bgTran := ebiten.GeoM{}
	bgTran.Translate(initialSeamOffsetX+bgOffsetX, initialSeamOffsetY+bgOffsetY)

	screen.DrawImage(bgLayer, &ebiten.DrawImageOptions{
		GeoM:          bgTran,
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
		ImageParts:    nil,
		Parts:         nil,
		SourceRect:    nil,
	})

	iconTran := ebiten.TranslateGeo(ScreenWidth/2 - 16, ScreenHeight-100)

	screen.DrawImage(playerIcon, &ebiten.DrawImageOptions{
		GeoM:          iconTran,
		ColorM:        ebiten.ColorM{},
		CompositeMode: 0,
		Filter:        0,
		ImageParts:    nil,
		Parts:         nil,
		SourceRect:    nil,
	})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", g.y), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", bgOffsetX), 100, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", bgOffsetY), 0, 10)
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

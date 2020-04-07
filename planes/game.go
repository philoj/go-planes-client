package planes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	_ "image/jpeg" // required for the image file loading to work. see ebitenutil.NewImageFromFile
	_ "image/png"
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

	initialVelocity = -0.5
	defaultRotation = 0.01
)

var bgLayer, playerIcon *ebiten.Image

func init() {
	fmt.Println("init")
	tile, _, _ := ebitenutil.NewImageFromFile("planes/assets/bg.jpg", ebiten.FilterDefault)
	player, _, _ := ebitenutil.NewImageFromFile("planes/assets/icon_orig.png", ebiten.FilterDefault)
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
	heading float64
	x       float64
	dx      float64
	y       float64
	dy      float64
}

func NewGame() *Game {
	return &Game{
		heading: math.Pi / 2,
		x:       0.0,
		y:       0.0,
	}
}

func (g *Game) Update(screen *ebiten.Image) error {

	if ebiten.IsKeyPressed(ebiten.KeyLeft) && !ebiten.IsKeyPressed(ebiten.KeyRight) {
		fmt.Println("left")
		g.rotate(-defaultRotation)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && !ebiten.IsKeyPressed(ebiten.KeyLeft) {
		fmt.Println("Right")
		g.rotate(+defaultRotation)
	}
	g.move(initialVelocity)
	return g.Draw(screen)
}
func (g *Game) Draw(screen *ebiten.Image) error {
	screen.DrawImage(bgLayer, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-16, -16)
	op.GeoM.Rotate(g.heading)
	op.GeoM.Translate((ScreenWidth/2)+g.x, ScreenHeight-100+g.y)
	screen.DrawImage(playerIcon, op)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", g.x), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", g.y), 100, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", degrees(g.heading)), 100, 10)
	return nil
}

func (g *Game) move(delta float64) {
	g.dx, g.dy = delta*math.Cos(g.heading), delta*math.Sin(g.heading)
	g.x+=g.dx
	g.y+=g.dy
	//bgOffsetX := math.Mod(g.x, seamOffsetLimitX)
	//bgOffsetY := math.Mod(g.y, seamOffsetLimitY)
	//bgTran := ebiten.GeoM{}
	//bgTran.Translate(initialSeamOffsetX+bgOffsetX, initialSeamOffsetY+bgOffsetY)

	//screen.DrawImage(bgLayer, &ebiten.DrawImageOptions{})

	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(-16, -16)
	//op.GeoM.Rotate(g.heading)
	//op.GeoM.Translate(ScreenWidth/2, ScreenHeight-100)
	//screen.DrawImage(playerIcon, op)
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", g.y), 0, 0)
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", degrees(g.heading)), 100, 10)
}

func (g *Game) rotate(dTheta float64) {
	g.heading += dTheta

}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func degrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

package planes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"
	"image/color"
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

	initialVelocity = 2.0
	defaultRotation = 0.03
	margin          = 50.0
	cameraWidth     = ScreenWidth - (2 * margin)
	cameraHeight    = ScreenHeight - (2 * margin)
	cameraRadius    = 100.0
	cameraVelocity  = 1.0

	debugEnabled = false
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
	input   string
	heading float64
	x       float64
	dx      float64
	y       float64
	dy      float64

	cameraX float64
	cameraY float64
}

func NewGame() *Game {
	return &Game{
		heading: -math.Pi / 2,
		x:       0.0,
		y:       0.0,
		cameraX: 0.0,
		cameraY: 0.0,
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && !ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.input = "left"
		g.rotate(-defaultRotation)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && !ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.input = "right"
		g.rotate(+defaultRotation)
	} else {
		g.input = ""
	}
	g.move(initialVelocity)
	return g.Draw(screen)
}
func (g *Game) Draw(screen *ebiten.Image) error {

	distanceToCamera := cartesianDistance(g.x, g.y, g.cameraX, g.cameraY)
	if distanceToCamera > cameraRadius {
		// player is beyond the edge, move camera
		g.cameraX += g.dx * cameraVelocity
		g.cameraY += g.dy * cameraVelocity
	}

	// bg
	bgOptions := &ebiten.DrawImageOptions{}

	bgOptions.GeoM.Translate(initialSeamOffsetX+math.Mod(-g.cameraX, seamOffsetLimitX), initialSeamOffsetY+math.Mod(-g.cameraY, seamOffsetLimitY))
	screen.DrawImage(bgLayer, bgOptions)

	// player
	playerOptions := &ebiten.DrawImageOptions{}
	playerOptions.GeoM.Translate(-16, -16)
	playerOptions.GeoM.Rotate(g.heading)
	playerOptions.GeoM.Translate((ScreenHeight/2)+g.x-g.cameraX, (ScreenWidth/2)+g.y-g.cameraY)
	screen.DrawImage(playerIcon, playerOptions)

	// debug info
	if debugEnabled {
		//drawRectangle(screen, margin, margin, cameraWidth, cameraHeight, colornames.Red)
		drawCircle(screen, ScreenWidth/2, ScreenHeight/2, cameraRadius, colornames.Red)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", g.x), 0, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", g.y), 100, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", degrees(g.heading)), 0, 10)
		ebitenutil.DebugPrintAt(screen, g.input, 100, 10)
	}
	return nil
}

func (g *Game) move(delta float64) {
	g.dx, g.dy = radialToXY(delta, g.heading)
	g.x += g.dx
	g.y += g.dy
	fmt.Println(g.dx, g.dy)
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

func drawRectangle(screen *ebiten.Image, x float64, y float64, width float64, height float64, c color.Color) {
	ebitenutil.DrawLine(screen, x, y, x+width, y, c)
	ebitenutil.DrawLine(screen, x+width, y, x+width, y+height, c)
	ebitenutil.DrawLine(screen, x+width, y+height, x, y+height, c)
	ebitenutil.DrawLine(screen, x, y+height, x, y, c)
}

func cartesianDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(((x2 - x1) * (x2 - x1)) + ((y2 - y1) * (y2 - y1)))
}

func radialToXY(radius, theta float64) (x, y float64) {
	return radius * math.Cos(theta), radius * math.Sin(theta)
}

const (
	segments = 16
	angle    = (2 * math.Pi) / segments
)

func drawCircle(screen *ebiten.Image, x, y, radius float64, c color.Color) {
	x1, y1 := x+radius, y
	for i := 0; i < segments; i++ {
		phi := float64(i) * angle
		phi2 := phi + angle
		x2 := x + (radius * math.Cos(phi2))
		y2 := y + (radius * math.Sin(phi2))
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, c)
		x1, y1 = x2, y2
	}
}

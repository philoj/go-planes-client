package planes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"
	_ "image/jpeg" // required for the image file loading to work. see ebitenutil.NewImageFromFile
	_ "image/png"
	"log"
	"math"
	"strconv"
	"strings"
	"watchYourSix/lobby"
)

const (
	ScreenWidth  = 600.0
	ScreenHeight = 600.0

	maxX               = ScreenWidth / 2
	maxY               = ScreenHeight / 2
	minX               = -maxX
	minY               = -maxY
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

	playerIconSize      = 64.0
	playerIconImageSize = 640.0
	playerIconScale     = playerIconSize / playerIconImageSize

	initialVelocity = 4.0
	defaultRotation = 0.03
	cameraRadius    = .3 * ScreenHeight
	cameraVelocity  = 1.0

	defaultHeading = -math.Pi / 2
	defaultX       = 0.0
	defaultY       = 0.0

	radarRadius = 120.0

	debugEnabled = false
)

var bgLayer, playerIcon *ebiten.Image

func init() {
	tile, _, _ := ebitenutil.NewImageFromFile("planes/assets/bg.jpg", ebiten.FilterDefault)
	player, _, _ := ebitenutil.NewImageFromFile("planes/assets/icon_orig.png", ebiten.FilterDefault)
	playerIcon, _ = ebiten.NewImage(playerIconSize, playerIconSize, ebiten.FilterDefault)
	pTran := ebiten.ScaleGeo(playerIconScale, playerIconScale)
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

func NewGame(playerId int) *Game {
	game := &Game{
		remotePlayers: make(map[int]*Player),
		cameraX:       0.0,
		cameraY:       0.0,
		Tick:          make(chan bool),
	}
	game.Player = NewPlayer(playerId, defaultHeading, defaultX, defaultY, game)
	go game.watchLobby()
	return game
}

type Game struct {
	*Player
	remotePlayers map[int]*Player

	Tick    chan bool
	input   string
	cameraX float64
	cameraY float64
}

func (g *Game) GetTicker() *chan bool {
	return &(g.Tick)
}

func (g *Game) Update(screen *ebiten.Image) error {
	// update Player
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && !ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.input = "left"
		g.Player.rotate(-defaultRotation)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && !ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.input = "right"
		g.Player.rotate(+defaultRotation)
	} else {
		g.input = ""
	}
	g.Player.move(initialVelocity)

	// broadcast location
	g.Tick <- true

	// draw
	return g.Draw(screen)
}
func (g *Game) Draw(screen *ebiten.Image) error {
	// draw ui
	distanceToCamera := cartesianDistance(g.Player.X, g.Player.Y, g.cameraX, g.cameraY)
	if distanceToCamera > cameraRadius {
		// Player is beyond the edge, move camera
		g.cameraX += g.Player.Vx * cameraVelocity
		g.cameraY += g.Player.Vy * cameraVelocity
	}

	// bg
	bgOptions := &ebiten.DrawImageOptions{}

	bgOptions.GeoM.Translate(initialSeamOffsetX+math.Mod(-g.cameraX, seamOffsetLimitX), initialSeamOffsetY+math.Mod(-g.cameraY, seamOffsetLimitY))
	screen.DrawImage(bgLayer, bgOptions)

	g.Player.Draw(screen)

	// draw other players
	for id := range g.remotePlayers {
		g.remotePlayers[id].Draw(screen)
	}

	// debug info
	if debugEnabled {
		//drawRectangle(screen, margin, margin, cameraWidth, cameraHeight, colornames.Red)
		drawCircle(screen, ScreenWidth/2, ScreenHeight/2, cameraRadius, colornames.Red)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f, %f", g.Player.X, g.Player.Y), 0, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", degrees(g.Player.Heading)), 0, 10)
		ebitenutil.DebugPrintAt(screen, g.input, 100, 10)
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) GetState() []byte {
	return []byte(fmt.Sprintf(
		"%d,%f,%f,%f,%f,%f",
		g.Player.Id,
		g.Player.X,
		g.Player.Y,
		g.Player.Vx,
		g.Player.Vy,
		g.Player.Heading,
	))
}

func (g *Game) updateRemotePlayer(dataByte []byte) {
	// Id, X, Y, Vx, Vy, Heading
	data := strings.Split(string(dataByte), ",")

	id, _ := strconv.Atoi(data[0])
	x, _ := strconv.ParseFloat(data[1], 64)
	y, _ := strconv.ParseFloat(data[2], 64)
	vx, _ := strconv.ParseFloat(data[3], 64)
	vy, _ := strconv.ParseFloat(data[4], 64)
	h, _ := strconv.ParseFloat(data[5], 64)
	player := g.remotePlayers[id]
	log.Println("player", player)
	if player == nil {
		// add new p without value for Game
		p := NewPlayer(id, h, x, y, g)
		p.isRemotePlayer = true
		p.Vx = vx
		p.Vy = vy
		p.Heading = h
		g.remotePlayers[id] = p
		log.Println("Added Player", p.Id)
	} else {
		// update existing player
		player.Update(x, y, vx, vy, h)
	}
}
func (g *Game) watchLobby() {
	for {
		select {
		case p := <-lobby.Lobby:
			g.updateRemotePlayer(p)
		}
	}
}

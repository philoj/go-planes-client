package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/rakyll/statik/fs"
	"goplanesclient/geometry"
	"goplanesclient/lobby"
	"goplanesclient/physics"
	"goplanesclient/players"
	"goplanesclient/plot"
	"goplanesclient/render"
	_ "goplanesclient/statik"
	"goplanesclient/touch"
	"image"
	_ "image/jpeg" // required for the image file loading to work. see ebitenutil.NewImageFromFile
	_ "image/png"
	"log"
	"strconv"
	"strings"
	"time"
)

func NewGame(playerId int, debug bool, host, path string) *Planes {
	game := &Planes{
		remotePlayers: make(map[int]*players.Player),
		tick:          make(chan bool, 1000),
		frameRate:     time.NewTicker(30 * time.Millisecond),
		debug:         debug,
		touch:         touch.NewTouchController(),
		images:        images,
	}
	game.player = players.NewPlayer(playerId, true, defaultX, defaultY, defaultHeading, 0, 0)
	go game.watchLobby()
	go lobby.JoinLobby(game, host, path)
	return game
}

type Planes struct {
	player        *players.Player
	remotePlayers map[int]*players.Player

	debug bool

	tick      chan bool
	input     string
	touch     touch.Controller
	frameRate *time.Ticker

	images        map[string]*imageInfo
	camera        *render.Camera
	cameraTracker physics.Tracker

	radarRadius                 float64
	outsideWidth, outsideHeight int
}

func (g *Planes) Update() error {
	// update Player
	g.input = ""
	g.touch.Read()

	if (ebiten.IsKeyPressed(ebiten.KeyLeft) && !ebiten.IsKeyPressed(ebiten.KeyRight)) ||
		g.touch.IsButtonPressed(leftTouchButtonId) && !g.touch.IsButtonPressed(rightTouchButtonId) {
		g.input = "left"
		g.player.Rotate(-defaultRotation)
	}
	if (ebiten.IsKeyPressed(ebiten.KeyRight) && !ebiten.IsKeyPressed(ebiten.KeyLeft)) ||
		g.touch.IsButtonPressed(rightTouchButtonId) && !g.touch.IsButtonPressed(leftTouchButtonId) {
		g.input = "right"
		g.player.Rotate(+defaultRotation)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) ||
		g.touch.IsButtonPressed(leftTouchButtonId) && g.touch.IsButtonPressed(rightTouchButtonId) {
		g.input += "+up"
		g.player.Move(defaultAcceleration * initialVelocity)
	}
	g.player.Move(initialVelocity)

	// broadcast location
	g.tick <- true

	// update camera location
	g.cameraTracker.UpdateFollower()
	return nil
}

func (g *Planes) Draw(screen *ebiten.Image) {
	// background
	bgTranslation := g.camera.Location().Negate() // negative of camera location
	plot.LaySquareTiledImage(screen, g.images[bgImageAssetId].image, bgTranslation, g.camera.Width, -1)

	// player
	g.camera.DrawObject(screen, g.images[iconImageAssetId].image, g.player.Mover)

	// draw other players
	for id := range g.remotePlayers {
		g.camera.DrawObject(screen, g.images[iconImageAssetId].image, g.remotePlayers[id].Mover)
	}

	// debug info
	if g.debug {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("players X: %f Y: %f H: %f", g.player.Location().I, g.player.Location().J, geometry.Degrees(g.player.Heading())), 0, 0)
		ebitenutil.DebugPrintAt(screen, g.input, 100, 10)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera X: %f Y: %f",
			g.camera.Mover.Location().I, g.camera.Mover.Location().J), 0, 50)
	}
}

func (g *Planes) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if outsideHeight != g.outsideHeight || outsideWidth != g.outsideWidth {
		// reload viewport
		g.outsideWidth, g.outsideHeight = outsideWidth, outsideHeight
		g.loadViewPort()
	}
	return outsideWidth, outsideHeight
}

func (g *Planes) GetState() []byte {
	return []byte(fmt.Sprintf(
		"%d,%f,%f,%f,%f,%f",
		g.player.Id,
		g.player.Location().I,
		g.player.Location().J,
		g.player.Velocity().I,
		g.player.Velocity().J,
		g.player.Heading(),
	))
}

func (g *Planes) GetTicker() *chan bool {
	return &(g.tick)
}

func (g *Planes) updateRemotePlayer(dataByte []byte) {
	// Id, X, Y, Vx, Vy, Heading
	data := strings.Split(string(dataByte), ",")

	id, _ := strconv.Atoi(data[0])
	x, _ := strconv.ParseFloat(data[1], 64)
	y, _ := strconv.ParseFloat(data[2], 64)
	vx, _ := strconv.ParseFloat(data[3], 64)
	vy, _ := strconv.ParseFloat(data[4], 64)
	h, _ := strconv.ParseFloat(data[5], 64)
	player := g.remotePlayers[id]
	if player == nil {
		// add new p without value for Game
		p := players.NewPlayer(id, false, x, y, h, vx, vy)
		g.remotePlayers[id] = p
		log.Println("Added Player", p.Id)
	} else {
		// update existing player
		player.Reset(x, y, vx, vy, h)
	}
}

func (g *Planes) watchLobby() {
	for {
		p := <-lobby.Lobby
		g.updateRemotePlayer(p)
	}
}

func (g *Planes) loadViewPort() {
	fWidth, fHeight := float64(g.outsideWidth), float64(g.outsideHeight)
	log.Print(fWidth, fHeight)
	g.camera = render.NewCamera(0, 0, 0, 0, 0, fWidth, fHeight)
	g.radarRadius = fHeight / 2

	// load images
	embeddedFs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	iconSize := g.outsideHeight / 10
	for imgId, imgInf := range g.images {
		log.Print(imgId)
		// Calculate image render sizes
		g.images[imgId].targetSize = imageSize{}
		switch imgId {
		case bgImageAssetId:
			size := g.outsideWidth * 3
			g.images[imgId].targetSize.width, g.images[imgId].targetSize.height = size, size
			break
		case iconImageAssetId:
			g.images[imgId].targetSize.width, g.images[imgId].targetSize.height = iconSize, iconSize
			break
		case blipImageAssetId:
			g.images[imgId].targetSize.width, g.images[imgId].targetSize.height = iconSize, iconSize
			break
		}

		// Open image files from embedded assets
		f, err := embeddedFs.Open(imgInf.path)
		if err != nil {
			log.Fatal(fmt.Errorf("fail to open %s: %s", imgInf.path, err))
		}
		//noinspection GoDeferInLoop
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal(fmt.Errorf("fail to close %s: %s", imgInf.path, err))
			}
		}()

		// Decode and create im memory ebiten images
		i, _, err := image.Decode(f)
		if err != nil {
			log.Fatal(fmt.Errorf("fail to decode %s: %s", imgInf.path, err))
		}
		original := ebiten.NewImageFromImage(i)
		g.images[imgId].image = ebiten.NewImage(imgInf.targetSize.width, imgInf.targetSize.height)
		transform := ebiten.GeoM{}
		transform.Scale(
			// Scale from original size to target size calculated earlier
			float64(imgInf.targetSize.width)/float64(imgInf.originalSize.width),
			float64(imgInf.targetSize.height)/float64(imgInf.originalSize.height))
		g.images[imgId].image.DrawImage(original, &ebiten.DrawImageOptions{
			GeoM: transform,
		})
	}

	// Tracker to make the camera follow player smoothly
	g.cameraTracker = physics.NewSimpleTracker(g.camera, g.player, fWidth/2, fWidth/2, cameraVelocity)

	// Mount all touch buttons on the touch controller
	for _, b := range allButtons(fWidth, fHeight) {
		g.touch.Mount(b)
	}
}

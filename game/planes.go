package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/rakyll/statik/fs"
	"goplanesclient/geometry"
	"goplanesclient/lobby"
	"goplanesclient/physics"
	"goplanesclient/players"
	"goplanesclient/plot"
	"goplanesclient/render"
	"goplanesclient/touch"
	"image"
	_ "image/jpeg" // required for the image file loading to work. see ebitenutil.NewImageFromFile
	_ "image/png"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"

	_ "goplanesclient/statik"
)

const (
	bgImageSize         = 5000.0
	playerIconImageSize = 640.0
	blipIconImageSize   = 320.0

	initialVelocity     = 4
	defaultAcceleration = 1
	defaultRotation     = 0.03
	cameraVelocity      = 0.1

	defaultHeading = -math.Pi / 2
	defaultX       = 0.0
	defaultY       = 0.0
)

const (
	BgImageAssetId   = "tile"
	IconImageAssetId = "players"
	BlipImageAssetId = "blip"

	leftTouchButtonId  = "left"
	rightTouchButtonId = "right"
)

func NewGame(playerId int, debug bool, host, path string) *Planes {
	game := &Planes{
		remotePlayers:   make(map[int]*players.Player),
		tick:            make(chan bool),
		viewPortLoading: sync.Once{},
		initComplete:    make(chan bool),
		debug:           debug,
		touch:           touch.NewTouchController(),
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

	tick  chan bool
	input string
	touch touch.Controller

	viewPortLoading sync.Once
	initComplete    chan bool
	images          map[string]*imageInfo
	camera          *render.Camera
	cameraTracker   physics.Tracker

	radarRadius float64
}
type imageSize struct {
	width, height int
}
type imageInfo struct {
	path                     string
	originalSize, targetSize imageSize
	image                    *ebiten.Image
}

func (g *Planes) Update(screen *ebiten.Image) error {
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

	// draw
	select {
	case <-g.initComplete:
	}
	// update camera location
	g.cameraTracker.UpdateFollower()
	return g.Draw(screen)
}
func (g *Planes) Draw(screen *ebiten.Image) error {
	// background
	bgTranslation := g.camera.Location().Vector().Negate() // negative of camera location
	plot.LaySquareTiledImage(screen, g.images[BgImageAssetId].image, bgTranslation, g.camera.Width, -1)

	// player
	g.camera.DrawObject(screen, g.images[IconImageAssetId].image, g.player.Mover)

	// draw other players
	for id := range g.remotePlayers {
		g.camera.DrawObject(screen, g.images[IconImageAssetId].image, g.remotePlayers[id].Mover)
	}

	// debug info
	if g.debug {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("players X: %f Y: %f H: %f", g.player.Location().X, g.player.Location().Y, geometry.Degrees(g.player.Heading())), 0, 0)
		ebitenutil.DebugPrintAt(screen, g.input, 100, 10)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera X: %f Y: %f",
			g.camera.Mover.Location().X, g.camera.Mover.Location().Y), 0, 50)
	}
	return nil
}

func (g *Planes) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// viewport is initialized in first call
	go g.viewPortLoading.Do(func() { g.loadViewPort(outsideWidth, outsideHeight) })
	return outsideWidth, outsideHeight
}

func (g *Planes) GetState() []byte {
	return []byte(fmt.Sprintf(
		"%d,%f,%f,%f,%f,%f",
		g.player.Id,
		g.player.Location().X,
		g.player.Location().Y,
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
		select {
		case p := <-lobby.Lobby:
			g.updateRemotePlayer(p)
		}
	}
}

func (g *Planes) loadViewPort(outsideWidth, outsideHeight int) {
	fWidth, fHeight := float64(outsideWidth), float64(outsideHeight)
	g.camera = render.NewCamera(0, 0, 0, 0, 0, fWidth, fHeight)
	g.radarRadius = fHeight / 2

	// load images
	iconSize := outsideHeight / 10
	g.images = map[string]*imageInfo{
		BgImageAssetId: {
			path: "/bg.jpg",
			originalSize: imageSize{
				width:  bgImageSize,
				height: bgImageSize,
			},
			targetSize: imageSize{
				width:  outsideWidth * 3,
				height: outsideWidth * 3,
			},
		},
		IconImageAssetId: {
			path: "/icon_orig.png",
			originalSize: imageSize{
				width:  playerIconImageSize,
				height: playerIconImageSize,
			},
			targetSize: imageSize{
				width: iconSize, height: iconSize,
			},
		},
		BlipImageAssetId: {
			path: "/blip.png",
			originalSize: imageSize{
				width:  blipIconImageSize,
				height: blipIconImageSize,
			},
			targetSize: imageSize{
				width: iconSize, height: iconSize,
			},
		},
	}
	embeddedFs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	for imgId, imgInf := range g.images {
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
		i, _, err := image.Decode(f)
		if err != nil {
			log.Fatal(fmt.Errorf("fail to decode %s: %s", imgInf.path, err))
		}
		original, _ := ebiten.NewImageFromImage(i, ebiten.FilterDefault)
		canvas, _ := ebiten.NewImage(imgInf.targetSize.width, imgInf.targetSize.height, ebiten.FilterDefault)
		transform := ebiten.GeoM{}
		transform.Scale(
			float64(imgInf.targetSize.width)/float64(imgInf.originalSize.width),
			float64(imgInf.targetSize.height)/float64(imgInf.originalSize.height))
		canvas.DrawImage(original, &ebiten.DrawImageOptions{
			GeoM: transform,
		})
		g.images[imgId].image = canvas
	}

	g.cameraTracker = physics.NewSimpleTracker(g.camera, g.player, fWidth/2, fWidth/2, cameraVelocity)

	// touch buttons
	buttons := []touch.Button{
		touch.NewButton(
			leftTouchButtonId, geometry.Point{
				X: 0,
				Y: 0,
			}, geometry.ClosedPolygon{
				geometry.Point{
					X: 0,
					Y: 0,
				},
				geometry.Point{
					X: fWidth / 2,
					Y: 0,
				},
				geometry.Point{
					X: fWidth / 2,
					Y: fHeight,
				},
				geometry.Point{
					X: 0,
					Y: fHeight,
				},
			}),
		touch.NewButton(
			rightTouchButtonId, geometry.Point{
				X: fWidth / 2,
				Y: 0,
			}, geometry.ClosedPolygon{
				geometry.Point{
					X: 0,
					Y: 0,
				},
				geometry.Point{
					X: fWidth / 2,
					Y: 0,
				},
				geometry.Point{
					X: fWidth / 2,
					Y: fHeight,
				},
				geometry.Point{
					X: 0,
					Y: fHeight,
				},
			}),
	}
	for _, b := range buttons {
		g.touch.Mount(b)
	}
	close(g.initComplete)
}

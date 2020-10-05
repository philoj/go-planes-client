package planes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/rakyll/statik/fs"
	"image"
	_ "image/jpeg" // required for the image file loading to work. see ebitenutil.NewImageFromFile
	_ "image/png"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"watchYourSix/lobby"

	_ "watchYourSix/statik"
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
	IconImageAssetId = "player"
	BlipImageAssetId = "blip"

	leftTouchButtonId  = "left"
	rightTouchButtonId = "right"
)

func NewGame(debug bool) *Planes {
	game := &Planes{
		remotePlayers: make(map[int]*Player),
		camera: &camera{
			PointObject: &PointObject{
				location: &Point{},
				velocity: &Vector{},
			},
		},
		Tick:            make(chan bool),
		viewPortLoading: sync.Once{},
		initComplete:    make(chan bool),
		debug:           debug,
		touch: &touchController{
			buttons: make(map[string]*buttonController),
			state:   make(map[string]bool),
		},
	}
	game.player = NewPlayer(0, defaultHeading, defaultX, defaultY)
	go game.watchLobby()
	go lobby.JoinLobby(game)
	return game
}

type Planes struct {
	player        *Player
	remotePlayers map[int]*Player

	debug bool

	Tick  chan bool
	input string
	touch *touchController

	viewPortLoading sync.Once
	initComplete    chan bool
	images          map[string]*imageInfo
	camera          *camera
	cameraTracker   TrackerInterface

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
	g.touch.ProcessTouch()

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
	g.Tick <- true

	// draw
	select {
	case <-g.initComplete:
	}
	// update camera location
	g.cameraTracker.UpdateTarget()
	return g.Draw(screen)
}
func (g *Planes) Draw(screen *ebiten.Image) error {
	// background
	bgTranslation := g.camera.location.Vector().Negate() // negative of camera location
	laySquareTiledImage(screen, g.images[BgImageAssetId].image, bgTranslation, g.camera.width, -1)

	// player
	g.camera.DrawObject(screen, g.images[IconImageAssetId].image, *g.player.PointObject)

	// draw other players
	for id := range g.remotePlayers {
		g.camera.DrawObject(screen, g.images[IconImageAssetId].image, *g.remotePlayers[id].PointObject)
	}

	// debug info
	if g.debug {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("player X: %f Y: %f H: %f", g.player.location.X, g.player.location.Y, degrees(g.player.heading)), 0, 0)
		ebitenutil.DebugPrintAt(screen, g.input, 100, 10)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera X: %f Y: %f",
			g.camera.PointObject.location.X, g.camera.PointObject.location.Y), 0, 50)
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
		g.player.location.X,
		g.player.location.Y,
		g.player.velocity.I,
		g.player.velocity.J,
		g.player.heading,
	))
}

func (g *Planes) GetTicker() *chan bool {
	return &(g.Tick)
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
		p := NewPlayer(id, h, x, y)
		p.isRemotePlayer = true
		p.velocity.I = vx
		p.velocity.J = vy
		p.heading = h
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
	g.camera.width = float64(outsideWidth)
	g.camera.height = float64(outsideHeight)
	g.radarRadius = g.camera.height / 2

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

	g.cameraTracker = NewSimpleTracker(g.camera, g.player, g.camera.width/2, g.camera.width/2, cameraVelocity)

	// touch buttons
	buttons := []touchButton{
		{
			id: leftTouchButtonId,
			location: Point{
				X: 0,
				Y: 0,
			},
			relativeGeometry: ClosedPolygon{
				Point{
					X: 0,
					Y: 0,
				},
				Point{
					X: g.camera.width / 2,
					Y: 0,
				},
				Point{
					X: g.camera.width / 2,
					Y: g.camera.height,
				},
				Point{
					X: 0,
					Y: g.camera.height,
				},
			},
		},
		{
			id: rightTouchButtonId,
			location: Point{
				X: g.camera.width / 2,
				Y: 0,
			},
			relativeGeometry: ClosedPolygon{
				Point{
					X: 0,
					Y: 0,
				},
				Point{
					X: g.camera.width / 2,
					Y: 0,
				},
				Point{
					X: g.camera.width / 2,
					Y: g.camera.height,
				},
				Point{
					X: 0,
					Y: g.camera.height,
				},
			},
		},
	}
	for _, b := range buttons {
		g.touch.Mount(&b)
	}
	close(g.initComplete)
}

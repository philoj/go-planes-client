package planes

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"
	"log"
	"math"
)

type Player struct {
	*Game
	isRemotePlayer bool
	Id             int
	Heading        float64
	X              float64
	Vx             float64
	Y              float64
	Vy             float64
}

func NewPlayer(id int, heading, X, Y float64, game *Game) *Player {
	return &Player{
		Game:    game,
		Id:      id,
		Heading: -math.Pi / 2,
		X:       0.0,
		Y:       0.0,
	}
}

func (p *Player) Update(x, y, vx, vy, heading float64) {
	log.Println("Updating player", x, y)
	p.X = x
	p.Y = y
	p.Vx = vx
	p.Vy = vy
	p.Heading = heading
}

func (p *Player) Draw(screen *ebiten.Image) {
	// draw Player
	x := p.X - p.Game.cameraX
	y := p.Y - p.Game.cameraY
	if !p.isRemotePlayer || (x > minX && x < maxX && y > minY && y < maxY) {
		playerOptions := &ebiten.DrawImageOptions{}
		playerOptions.GeoM.Translate(-playerIconSize, -playerIconSize)
		playerOptions.GeoM.Rotate(p.Heading)
		playerOptions.GeoM.Translate(maxX+x, maxY+y)
		screen.DrawImage(playerIcon, playerOptions)
	} else if p.isRemotePlayer {
		// remote player outside camera. draw a pointer
		theta := math.Atan(y / x)
		blipX, blipY := RadialToXY(radarRadius, theta)

		// adjust for the sign inversions possible in trigonometry
		if (x < 0 && blipX > 0) || (x > 0 && blipX < 0) {
			blipX = -blipX
		}
		if (y < 0 && blipY > 0) || (y > 0 && blipY < 0) {
			blipY = -blipY
		}
		drawCircle(screen, maxX+blipX, maxY+blipY, 20, colornames.Black)
		ebitenutil.DrawLine(screen, maxX, maxY, maxX+x, maxY+y, colornames.Black)
	}
}

func (p *Player) move(delta float64) {
	p.Vx, p.Vy = RadialToXY(delta, p.Heading)
	p.X += p.Vx
	p.Y += p.Vy
}

func (p *Player) rotate(dTheta float64) {
	p.Heading += dTheta
	p.Heading = math.Mod(p.Heading, 2*math.Pi)
}

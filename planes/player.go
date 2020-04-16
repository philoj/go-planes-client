package planes

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"
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
	//log.Println("Updating player", x, y)
	p.X = x
	p.Y = y
	p.Vx = vx
	p.Vy = vy
	p.Heading = heading
}

func (p *Player) Draw(screen *ebiten.Image) {
	// draw Player
	lx := p.X - p.Game.cameraX
	ly := p.Y - p.Game.cameraY
	if !p.isRemotePlayer || (lx > minX && lx < maxX && ly > minY && ly < maxY) {
		playerOptions := &ebiten.DrawImageOptions{}
		playerOptions.GeoM.Translate(-playerIconSize/2, -playerIconSize/2)
		playerOptions.GeoM.Rotate(p.Heading)
		playerOptions.GeoM.Translate(maxX+lx, maxY+ly)
		screen.DrawImage(playerIcon, playerOptions)
	} else if p.isRemotePlayer {

		blipX, blipY := BisectLine(0, 0, lx, ly, radarRadius)
		blipTran := ebiten.TranslateGeo(-blipIconSize/2, -blipIconSize/2)
		blipTran.Rotate(Theta(blipX, blipY))
		blipTran.Translate(maxX+blipX, maxY+blipY)
		screen.DrawImage(radarBlip, &ebiten.DrawImageOptions{
			GeoM: blipTran,
		})
		if debugEnabled {
			drawCircle(screen, maxX+blipX, maxY+blipY, 20, colornames.Black)
			ebitenutil.DrawLine(screen, maxX, maxY, maxX+lx, maxY+ly, colornames.Black)
		}
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

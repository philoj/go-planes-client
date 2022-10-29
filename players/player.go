package players

import (
	"goplanesclient/physics"
)

type Player struct {
	physics.Mover
	isSelf bool
	Id     int
}

func NewPlayer(id int, isSelf bool, x, y, theta, i, j float64) *Player {
	return &Player{
		Mover:  physics.NewMovingObject(x, y, i, j, theta),
		isSelf: isSelf,
		Id:     id,
	}
}

func (p *Player) Reset(x, y, i, j, theta float64) {
	p.Mover = physics.NewMovingObject(x, y, i, j, theta)
}

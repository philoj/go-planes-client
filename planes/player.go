package planes

type Player struct {
	*PointObject
	isRemotePlayer bool
	Id             int
}

func NewPlayer(id int, heading, X, Y float64) *Player {
	return &Player{
		PointObject: &PointObject{
			location: &Point{
				X: X,
				Y: Y,
			},
			velocity: &Vector{},
			heading:  heading,
		},
		Id:   id,
	}
}

func (p *Player) Reset(x, y, vx, vy, heading float64) {
	p.location = &Point{X: x, Y: y}
	p.velocity = &Vector{I: vx, J: vy}
	p.heading = heading
}

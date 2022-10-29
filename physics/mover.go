package physics

import (
	"goplanesclient/geometry"
	"math"
)

type Mover interface {
	Location() geometry.Vector
	Velocity() geometry.Vector
	Heading() float64
	Move(delta float64)
	Rotate(dTheta float64)
	Turn(heading float64)
	Jump(location geometry.Vector)
}

func NewMover(x, y, i, j, theta float64) Mover {
	return &movingObject{
		geometry.Vector{I: x, J: y}, geometry.Vector{I: i, J: j}, theta,
	}
}

type movingObject struct {
	location geometry.Vector
	velocity geometry.Vector
	heading  float64 // radians
}

func (p *movingObject) Location() geometry.Vector {
	return p.location
}
func (p *movingObject) Velocity() geometry.Vector {
	return p.velocity
}
func (p *movingObject) Heading() float64 {
	return p.heading
}
func (p *movingObject) Move(delta float64) {
	p.velocity.I, p.velocity.J = geometry.RadialToXY(delta, p.heading)
	p.location.I += p.velocity.I
	p.location.J += p.velocity.J
}

func (p *movingObject) Rotate(dTheta float64) {
	p.heading += dTheta
	p.heading = math.Mod(p.heading, 2*math.Pi)
}
func (p *movingObject) Turn(heading float64) {
	p.heading = heading
}
func (p *movingObject) Jump(location geometry.Vector) {
	p.location.I, p.location.J = location.I, location.J
}

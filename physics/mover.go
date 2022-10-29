package physics

import (
	"goplanesclient/geometry"
	"math"
)

type Mover interface {
	Location() geometry.Point
	Velocity() geometry.Vector
	Heading() float64
	Move(delta float64)
	Rotate(dTheta float64)
	Turn(heading float64)
	Jump(location geometry.Point)
}

func NewMover(x, y, i, j, theta float64) Mover {
	return &movingObject{
		&geometry.Point{X: x, Y: y}, &geometry.Vector{I: i, J: j}, theta,
	}
}

type movingObject struct {
	location *geometry.Point
	velocity *geometry.Vector
	heading  float64 // radians
}

func (p *movingObject) Location() geometry.Point {
	return *p.location
}
func (p *movingObject) Velocity() geometry.Vector {
	return *p.velocity
}
func (p *movingObject) Heading() float64 {
	return p.heading
}
func (p *movingObject) Move(delta float64) {
	p.velocity.I, p.velocity.J = geometry.RadialToXY(delta, p.heading)
	p.location.X += p.velocity.I
	p.location.Y += p.velocity.J
}

func (p *movingObject) Rotate(dTheta float64) {
	p.heading += dTheta
	p.heading = math.Mod(p.heading, 2*math.Pi)
}
func (p *movingObject) Turn(heading float64) {
	p.heading = heading
}
func (p *movingObject) Jump(location geometry.Point) {
	p.location.X, p.location.Y = location.X, location.Y
}

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

func NewMovingObject(x, y, i, j, theta float64) *MovingObject {
	return &MovingObject{
		&geometry.Point{X: x, Y: y}, &geometry.Vector{I: i, J: j}, theta,
	}
}

type MovingObject struct {
	location *geometry.Point
	velocity *geometry.Vector
	heading  float64 // radians
}

func (p *MovingObject) Location() geometry.Point {
	return *p.location
}
func (p *MovingObject) Velocity() geometry.Vector {
	return *p.velocity
}
func (p *MovingObject) Heading() float64 {
	return p.heading
}
func (p *MovingObject) Move(delta float64) {
	p.velocity.I, p.velocity.J = geometry.RadialToXY(delta, p.heading)
	p.location.X += p.velocity.I
	p.location.Y += p.velocity.J
}

func (p *MovingObject) Rotate(dTheta float64) {
	p.heading += dTheta
	p.heading = math.Mod(p.heading, 2*math.Pi)
}
func (p *MovingObject) Turn(heading float64) {
	p.heading = heading
}
func (p *MovingObject) Jump(location geometry.Point) {
	p.location.X, p.location.Y = location.X, location.Y
}

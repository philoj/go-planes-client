package planes

import "math"

type PointObjectInterface interface {
	Location() Point
	Velocity() Vector
	Heading() float64
	Move(delta float64)
	Rotate(dTheta float64)
	Turn(heading float64)
	Jump(location Point)
}

type PointObject struct {
	location *Point
	velocity *Vector
	heading  float64 // radians
}

func (p *PointObject) Location() Point {
	return *p.location
}
func (p *PointObject) Velocity() Vector {
	return *p.velocity
}
func (p *PointObject) Heading() float64 {
	return p.heading
}
func (p *PointObject) Move(delta float64) {
	p.velocity.I, p.velocity.J = RadialToXY(delta, p.heading)
	p.location.X += p.velocity.I
	p.location.Y += p.velocity.J
}

func (p *PointObject) Rotate(dTheta float64) {
	p.heading += dTheta
	p.heading = math.Mod(p.heading, 2*math.Pi)
}
func (p *PointObject) Turn(heading float64) {
	p.heading = heading
}
func (p *PointObject) Jump(location Point) {
	p.location.X, p.location.Y = location.X, location.Y
}

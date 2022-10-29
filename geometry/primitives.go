package geometry

import (
	"math"
)

type Point struct {
	X, Y float64
}

func (p Point) Vector() Vector {
	return Vector{
		I: p.X,
		J: p.Y,
	}
}

type Vector struct {
	I, J float64
}

func (v Vector) Negate() Vector {
	return Vector{
		I: -v.I,
		J: -v.J,
	}
}
func (v Vector) Add(v1 Vector) Vector {
	return Vector{
		I: v.I + v1.I,
		J: v.J + v1.J,
	}
}
func (v Vector) Size() float64 {
	return math.Sqrt((v.I * v.I) + (v.J * v.J))
}
func (v Vector) Point() Point {
	return Point{X: v.I, Y: v.J}
}

type Rectangle struct {
	Width, Height float64
}

type ClosedCurve interface {
	Inside(p Point) bool
}

type ClosedPolygon []Point

func (pg ClosedPolygon) Inside(pt Point) bool {
	if len(pg) < 3 {
		return false
	}
	in := rayIntersectsSegment(pt, pg[len(pg)-1], pg[0])
	for i := 1; i < len(pg); i++ {
		if rayIntersectsSegment(pt, pg[i-1], pg[i]) {
			in = !in
		}
	}
	return in
}

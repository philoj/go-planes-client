package geometry

import (
	"math"
)

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

type Rectangle struct {
	Width, Height float64
}

type ClosedCurve interface {
	Inside(p Vector) bool
}

type ClosedPolygon []Vector

func (pg ClosedPolygon) Inside(pt Vector) bool {
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

package geometry

import (
	"fmt"
	"math"
)

func CartesianDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(((x2 - x1) * (x2 - x1)) + ((y2 - y1) * (y2 - y1)))
}

func RadialToXY(radius, theta float64) (x, y float64) {
	return radius * math.Cos(theta), radius * math.Sin(theta)
}

/*
BisectLine
Return point (x,y) which bisects the line (x1,y1)-(x2,y2) with distance l from (x1,y1)
*/
func BisectLine(p1, p2 Point, d float64) Point {
	l := AxialDistance(p1, p2)
	theta := math.Atan(l.J / l.I)
	dx, dy := RadialToXY(d, theta)

	// adjust for the correct trigonometric quadrant
	if (l.I < 0 && dx > 0) || (l.I > 0 && dx < 0) {
		dx = -dx
	}
	if (l.J < 0 && dy > 0) || (l.J > 0 && dy < 0) {
		dy = -dy
	}
	return Point{X: p1.X + dx, Y: p1.Y + dy}
}

func BisectRectangle(p1, p2, rectMin, rectMax Point) Point {
	p := Point{
		X: p2.X,
		Y: p2.Y,
	}
	if rectMin.X > rectMax.X || rectMin.Y > rectMax.Y {
		panic(fmt.Errorf("invalid values for min and max %v, %v", rectMin, rectMax))
	}
	if p1.X > rectMin.X && p1.X < rectMax.X && p1.Y > rectMin.Y && p1.Y < rectMax.Y {
		if p2.X > rectMax.X {
			p.X = rectMax.X
		} else if p2.X < rectMin.X {
			p.X = rectMin.X
		}
		if p2.Y > rectMax.Y {
			p.Y = rectMax.Y
		} else if p2.Y < rectMin.Y {
			p.Y = rectMin.Y
		}
		if p.X == p2.X && p.Y == p2.Y {
			panic("invalid value for p2")
		}
		return p
	}
	panic("invalid value for p1")
}

func AxialDistance(p1, p2 Point) Vector {
	return Vector{I: p2.X - p1.X, J: p2.Y - p1.Y}
}

func Theta(v Vector) float64 {
	tan := v.J / v.I
	if tan == 0 {
		if v.I > 0 {
			return 0
		} else {
			return math.Pi
		}
	} else if tan < 0 {
		theta := math.Atan(-tan)
		if v.J > 0 {
			return math.Pi - theta
		} else {
			return -theta
		}
	} else {
		theta := math.Atan(tan)
		if v.J < 0 {
			return math.Pi + theta
		} else {
			return theta
		}
	}
}

func Degrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

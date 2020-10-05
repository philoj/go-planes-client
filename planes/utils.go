package planes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
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

func rayIntersectsSegment(p, a, b Point) bool {
	return (a.Y > p.Y) != (b.Y > p.Y) &&
		p.X < (b.X-a.X)*(p.Y-a.Y)/(b.Y-a.Y)+a.X
}

func CartesianDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(((x2 - x1) * (x2 - x1)) + ((y2 - y1) * (y2 - y1)))
}

func RadialToXY(radius, theta float64) (x, y float64) {
	return radius * math.Cos(theta), radius * math.Sin(theta)
}

/*
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

const (
	segments = 16
	angle    = (2 * math.Pi) / segments
)

func drawCircle(screen *ebiten.Image, x, y, radius float64, c color.Color) {
	x1, y1 := x+radius, y
	for i := 0; i < segments; i++ {
		phi := float64(i) * angle
		phi2 := phi + angle
		x2 := x + (radius * math.Cos(phi2))
		y2 := y + (radius * math.Sin(phi2))
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, c)
		x1, y1 = x2, y2
	}
}

func degrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

func drawRectangle(screen *ebiten.Image, x float64, y float64, width float64, height float64, c color.Color) {
	ebitenutil.DrawLine(screen, x, y, x+width, y, c)
	ebitenutil.DrawLine(screen, x+width, y, x+width, y+height, c)
	ebitenutil.DrawLine(screen, x+width, y+height, x, y+height, c)
	ebitenutil.DrawLine(screen, x, y+height, x, y, c)
}
func drawImage(screen, img *ebiten.Image, translate Vector, heading float64) {
	w, h := img.Size()
	rotScale := &ebiten.DrawImageOptions{}
	rotScale.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	rotScale.GeoM.Rotate(heading)
	rotScale.GeoM.Translate(translate.I, translate.J)
	screen.DrawImage(img, rotScale)
}
func laySquareTiledImage(screen, tile *ebiten.Image, originalTranslation Vector, tileSize float64, offsetCount int) {
	bgOptions := &ebiten.DrawImageOptions{}
	oc := float64(offsetCount)
	equivalentTranslation := Vector{
		I: math.Mod(originalTranslation.I, tileSize),
		J: math.Mod(originalTranslation.J, tileSize),
	}.Add(Vector{
		I: tileSize * oc,
		J: tileSize * oc,
	})
	bgOptions.GeoM.Translate(equivalentTranslation.I, equivalentTranslation.J)
	screen.DrawImage(tile, bgOptions)
}

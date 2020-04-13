package planes

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"math"
)

func cartesianDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(((x2 - x1) * (x2 - x1)) + ((y2 - y1) * (y2 - y1)))
}

func RadialToXY(radius, theta float64) (x, y float64) {
	return radius * math.Cos(theta), radius * math.Sin(theta)
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


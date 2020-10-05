package planes

import (
	"github.com/hajimehoshi/ebiten"
)

type camera struct {
	*PointObject
	width  float64
	height float64
}

// todo use top, bottom, etc in ui context
func (c *camera) LeftBoundary() float64 {
	return c.PointObject.location.X - (c.width / 2)
}
func (c *camera) RightBoundary() float64 {
	return c.PointObject.location.X + (c.width / 2)
}
func (c *camera) BottomBoundary() float64 {
	return c.PointObject.location.Y - (c.height / 2)
}
func (c *camera) TopBoundary() float64 {
	return c.PointObject.location.Y + (c.height / 2)
}
func (c *camera) Origin() Point {
	return Point{
		X: c.LeftBoundary(),
		Y: c.BottomBoundary(),
	}
}

func (c *camera) DrawObject(screen, img *ebiten.Image, p PointObject) {
	if p.location.X > c.LeftBoundary() && p.location.X < c.RightBoundary() && p.location.Y > c.BottomBoundary() && p.location.Y < c.TopBoundary() {
		drawImage(screen, img, AxialDistance(c.Origin(), *p.location), p.heading)
	}
}

package render

import (
	"github.com/hajimehoshi/ebiten"
	"goplanesclient/geometry"
	"goplanesclient/physics"
	"goplanesclient/plot"
)

func NewCamera(x, y, i, j, theta, w, h float64) *Camera {
	return &Camera{
		MovingObject: physics.NewMovingObject(x, y, i, j, theta),
		Rectangle: geometry.Rectangle{
			Width:  w,
			Height: h,
		},
	}
}

type Camera struct {
	geometry.Rectangle
	*physics.MovingObject
}

// todo use top, bottom, etc in ui context
func (c *Camera) LeftBoundary() float64 {
	return c.MovingObject.Location().X - (c.Width / 2)
}
func (c *Camera) RightBoundary() float64 {
	return c.MovingObject.Location().X + (c.Width / 2)
}
func (c *Camera) BottomBoundary() float64 {
	return c.MovingObject.Location().Y - (c.Height / 2)
}
func (c *Camera) TopBoundary() float64 {
	return c.MovingObject.Location().Y + (c.Height / 2)
}
func (c *Camera) Origin() geometry.Point {
	return geometry.Point{
		X: c.LeftBoundary(),
		Y: c.BottomBoundary(),
	}
}

func (c *Camera) DrawObject(screen, img *ebiten.Image, p physics.Mover) {
	if p.Location().X > c.LeftBoundary() && p.Location().X < c.RightBoundary() && p.Location().Y > c.BottomBoundary() && p.Location().Y < c.TopBoundary() {
		plot.DrawImage(screen, img, geometry.AxialDistance(c.Origin(), p.Location()), p.Heading())
	}
}

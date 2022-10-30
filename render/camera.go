package render

import (
	"github.com/hajimehoshi/ebiten/v2"
	"goplanesclient/geometry"
	"goplanesclient/physics"
	"goplanesclient/plot"
)

func NewCamera(x, y, i, j, theta, w, h float64) *Camera {
	return &Camera{
		Mover: physics.NewMover(x, y, i, j, theta),
		Rectangle: geometry.Rectangle{
			Width:  w,
			Height: h,
		},
	}
}

type Camera struct {
	geometry.Rectangle
	physics.Mover
}

// todo use top, bottom, etc in ui context
func (c *Camera) LeftBoundary() float64 {
	return c.Mover.Location().I - (c.Width / 2)
}
func (c *Camera) RightBoundary() float64 {
	return c.Mover.Location().I + (c.Width / 2)
}
func (c *Camera) BottomBoundary() float64 {
	return c.Mover.Location().J - (c.Height / 2)
}
func (c *Camera) TopBoundary() float64 {
	return c.Mover.Location().J + (c.Height / 2)
}
func (c *Camera) Origin() geometry.Vector {
	return geometry.Vector{
		I: c.LeftBoundary(),
		J: c.BottomBoundary(),
	}
}

func (c *Camera) DrawObject(screen, img *ebiten.Image, p physics.Mover) {
	if p.Location().I > c.LeftBoundary() && p.Location().I < c.RightBoundary() && p.Location().J > c.BottomBoundary() && p.Location().J < c.TopBoundary() {
		plot.DrawImage(screen, img, geometry.AxialDistance(c.Origin(), p.Location()), p.Heading())
	}
}

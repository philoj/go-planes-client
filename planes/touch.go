package planes

import (
	"github.com/hajimehoshi/ebiten"
	"time"
)

const refreshInterval = 200 * time.Millisecond

type touchButton struct {
	id               string
	location         Point
	relativeGeometry ClosedPolygon
}
type buttonController struct {
	button           *touchButton
	absoluteGeometry ClosedPolygon
}

type touchController struct {
	buttons  map[string]*buttonController
	state    map[string]bool
}

func (c *touchController) Mount(b *touchButton) {
	c.buttons[b.id] = &buttonController{
		button:           b,
		absoluteGeometry: nil,
	}
	for _, p := range b.relativeGeometry {
		c.buttons[b.id].absoluteGeometry = append(c.buttons[b.id].absoluteGeometry, p.Vector().Add(b.location.Vector()).Point())
	}
}

func (c *touchController) Locate(p Point) string {
	for id, b := range c.buttons {
		if b.absoluteGeometry.Inside(p) {
			return id
		}
	}
	return ""
}

func (c *touchController) ProcessTouch() {
	touchedButtonIds := make(map[string]bool)
	for _, tid := range ebiten.TouchIDs() {
		x, y := ebiten.TouchPosition(tid)
		if x != 0 && y != 0 {
			// todo save this conversion?
			p := Point{X: float64(x), Y: float64(y)}
			for id, b := range c.buttons {
				if b.absoluteGeometry.Inside(p) {
					touchedButtonIds[id] = true
					break
				}
			}
		}
	}
	c.state = touchedButtonIds
}
func (c *touchController) IsButtonPressed(id string) bool {
	pressed, ok := c.state[id]
	return ok && pressed
}

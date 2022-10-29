package touch

import (
	"github.com/hajimehoshi/ebiten"
	"goplanesclient/geometry"
	"time"
)

type Controller interface {
	Mount(b Button)
	Locate(p geometry.Point) string
	Read()
	IsButtonPressed(id string) bool
}

func NewTouchController() Controller {
	return &touchController{
		buttons: make(map[string]Button),
		state:   make(map[string]bool),
	}
}

const refreshInterval = 200 * time.Millisecond

type touchController struct {
	buttons map[string]Button
	state   map[string]bool
}

func (c *touchController) Mount(b Button) {
	c.buttons[b.Id()] = b
}

func (c *touchController) Locate(p geometry.Point) string {
	for id, b := range c.buttons {
		if b.Inside(p) {
			return id
		}
	}
	return ""
}

func (c *touchController) Read() {
	touchedIds := make(map[string]bool)
	for _, tid := range ebiten.TouchIDs() {
		x, y := ebiten.TouchPosition(tid)
		if x != 0 && y != 0 {
			// todo save this conversion?
			p := geometry.Point{X: float64(x), Y: float64(y)}
			for id, b := range c.buttons {
				if b.Inside(p) {
					touchedIds[id] = true
					break
				}
			}
		}
	}
	c.state = touchedIds
}
func (c *touchController) IsButtonPressed(id string) bool {
	pressed, ok := c.state[id]
	return ok && pressed
}

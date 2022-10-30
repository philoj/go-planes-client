package touch

import (
	"github.com/hajimehoshi/ebiten/v2"
	"goplanesclient/geometry"
)

type Controller interface {
	Mount(b Button)
	Locate(p geometry.Vector) string
	Read()
	IsButtonPressed(id string) bool
}

func NewTouchController() Controller {
	return &touchController{
		buttons: make(map[string]Button),
		state:   make(map[string]bool),
	}
}

type touchController struct {
	buttons map[string]Button
	state   map[string]bool
}

func (c *touchController) Mount(b Button) {
	c.buttons[b.Id()] = b
}

func (c *touchController) Locate(p geometry.Vector) string {
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
			p := geometry.Vector{I: float64(x), J: float64(y)}
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

package planes

import (
	"log"
	"math"
)

type TrackerInterface interface {
	UpdateTarget(delta float64)
}

func NewSimpleTracker(follower, leader PointObjectInterface, width, height, velocity float64) TrackerInterface {
	return &SimpleTracker{
		follower: follower,
		leader:   leader,
		maxX:     width / 2,
		maxY:     height / 2,
		velocity: velocity,
	}
}

type SimpleTracker struct {
	follower PointObjectInterface
	leader   PointObjectInterface
	maxX     float64
	maxY     float64
	velocity float64
}

func (t *SimpleTracker) UpdateTarget(delta float64) {
	d := AxialDistance(t.follower.Location(), t.leader.Location())
	if math.Abs(d.I) > t.maxX || math.Abs(d.J) > t.maxY {
		b := BisectRectangle(t.follower.Location(), t.leader.Location(), Point{
			X: t.follower.Location().X - t.maxX,
			Y: t.follower.Location().Y - t.maxY,
		}, Point{
			X: t.follower.Location().X + t.maxX,
			Y: t.follower.Location().Y + t.maxY,
		})
		v := AxialDistance(b, t.leader.Location())
		h := Theta(v)
		log.Println("tracker ",v, h)
		t.follower.Turn(h)
		t.follower.Move(v.Size() * t.velocity)
	}
}

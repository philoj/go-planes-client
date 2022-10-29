package physics

import (
	"goplanesclient/geometry"
	"math"
)

type Tracker interface {
	UpdateFollower()
}

func NewSimpleTracker(follower, leader Mover, width, height, velocity float64) Tracker {
	return &SimpleTracker{
		follower: follower,
		leader:   leader,
		maxX:     width / 2,
		maxY:     height / 2,
		velocity: velocity,
	}
}

type SimpleTracker struct {
	follower Mover
	leader   Mover
	maxX     float64
	maxY     float64
	velocity float64
}

func (t *SimpleTracker) UpdateFollower() {
	d := geometry.AxialDistance(t.follower.Location(), t.leader.Location())
	if math.Abs(d.I) > t.maxX || math.Abs(d.J) > t.maxY {
		b := geometry.BisectRectangle(t.follower.Location(), t.leader.Location(), geometry.Vector{
			I: t.follower.Location().I - t.maxX,
			J: t.follower.Location().J - t.maxY,
		}, geometry.Vector{
			I: t.follower.Location().I + t.maxX,
			J: t.follower.Location().J + t.maxY,
		})
		v := geometry.AxialDistance(b, t.leader.Location())
		h := geometry.Theta(v)
		t.follower.Turn(h)
		t.follower.Move(v.Size() * t.velocity)
	}
}

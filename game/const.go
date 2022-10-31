package game

import "math"

const (
	initialVelocity     = 4
	defaultAcceleration = 1
	defaultRotation     = 0.03
	cameraVelocity      = 0.1

	defaultHeading = -math.Pi / 2
	defaultX       = 0.0
	defaultY       = 0.0
)

const (
	bgImageAssetId   = "tile"
	iconImageAssetId = "players"
	blipImageAssetId = "blip"

	leftTouchButtonId  = "left"
	rightTouchButtonId = "right"
)

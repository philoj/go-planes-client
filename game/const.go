package game

import "math"

const (
	bgImageSize         = 5000.0
	playerIconImageSize = 640.0
	blipIconImageSize   = 320.0

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

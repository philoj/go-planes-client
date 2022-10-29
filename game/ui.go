package game

import (
	"goplanesclient/geometry"
	"goplanesclient/touch"
)

func allButtons(width, height float64) []touch.Button {
	return []touch.Button{
		touch.NewButton(
			leftTouchButtonId, geometry.Vector{
				I: 0,
				J: 0,
			}, geometry.ClosedPolygon{
				geometry.Vector{
					I: 0,
					J: 0,
				},
				geometry.Vector{
					I: width / 2,
					J: 0,
				},
				geometry.Vector{
					I: width / 2,
					J: height,
				},
				geometry.Vector{
					I: 0,
					J: height,
				},
			}),
		touch.NewButton(
			rightTouchButtonId, geometry.Vector{
				I: width / 2,
				J: 0,
			}, geometry.ClosedPolygon{
				geometry.Vector{
					I: 0,
					J: 0,
				},
				geometry.Vector{
					I: width / 2,
					J: 0,
				},
				geometry.Vector{
					I: width / 2,
					J: height,
				},
				geometry.Vector{
					I: 0,
					J: height,
				},
			}),
	}
}

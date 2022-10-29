package game

import (
	"goplanesclient/geometry"
	"goplanesclient/touch"
)

func allButtons(width, height float64) []touch.Button {
	return []touch.Button{
		touch.NewButton(
			leftTouchButtonId, geometry.Point{
				X: 0,
				Y: 0,
			}, geometry.ClosedPolygon{
				geometry.Point{
					X: 0,
					Y: 0,
				},
				geometry.Point{
					X: width / 2,
					Y: 0,
				},
				geometry.Point{
					X: width / 2,
					Y: height,
				},
				geometry.Point{
					X: 0,
					Y: height,
				},
			}),
		touch.NewButton(
			rightTouchButtonId, geometry.Point{
				X: width / 2,
				Y: 0,
			}, geometry.ClosedPolygon{
				geometry.Point{
					X: 0,
					Y: 0,
				},
				geometry.Point{
					X: width / 2,
					Y: 0,
				},
				geometry.Point{
					X: width / 2,
					Y: height,
				},
				geometry.Point{
					X: 0,
					Y: height,
				},
			}),
	}
}

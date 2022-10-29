package touch

import "goplanesclient/geometry"

type Button interface {
	geometry.ClosedCurve
	Id() string
	Location() geometry.Point
	Shape() geometry.ClosedPolygon
}

func NewButton(id string, location geometry.Point, shape geometry.ClosedPolygon) Button {
	absShape := make(geometry.ClosedPolygon, len(shape))
	for i, p := range shape {
		absShape[i] = p.Vector().Add(location.Vector()).Point()
	}
	return &touchButton{
		ClosedPolygon: absShape,
		id:            id,
		location:      location,
		shape:         shape,
	}
}

type touchButton struct {
	geometry.ClosedPolygon
	id       string
	location geometry.Point
	shape    geometry.ClosedPolygon
}

func (b *touchButton) Id() string {
	return b.id
}

func (b *touchButton) Location() geometry.Point {
	return b.location
}

func (b *touchButton) Shape() geometry.ClosedPolygon {
	return b.shape
}

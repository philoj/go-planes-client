package touch

import "goplanesclient/geometry"

type Button interface {
	geometry.ClosedCurve
	Id() string
	Location() geometry.Vector
	Shape() geometry.ClosedPolygon
}

func NewButton(id string, location geometry.Vector, shape geometry.ClosedPolygon) Button {
	absShape := make(geometry.ClosedPolygon, len(shape))
	for i, p := range shape {
		absShape[i] = p.Add(location)
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
	location geometry.Vector
	shape    geometry.ClosedPolygon
}

func (b *touchButton) Id() string {
	return b.id
}

func (b *touchButton) Location() geometry.Vector {
	return b.location
}

func (b *touchButton) Shape() geometry.ClosedPolygon {
	return b.shape
}

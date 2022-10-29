package geometry

func rayIntersectsSegment(p, a, b Point) bool {
	return (a.Y > p.Y) != (b.Y > p.Y) &&
		p.X < (b.X-a.X)*(p.Y-a.Y)/(b.Y-a.Y)+a.X
}

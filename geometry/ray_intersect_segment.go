package geometry

func rayIntersectsSegment(p, a, b Vector) bool {
	return (a.J > p.J) != (b.J > p.J) &&
		p.I < (b.I-a.I)*(p.J-a.J)/(b.J-a.J)+a.I
}

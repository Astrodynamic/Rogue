package domain

type Rect struct {
	Point
	W int
	H int
}

func NewRect(x, y, w, h int) Rect {
	return Rect{Point: Point{X: x, Y: y}, W: w, H: h}
}

func (r *Rect) Contains(p Point) bool {
	return p.X >= r.Point.X && p.X < r.Point.X+r.W && p.Y >= r.Point.Y && p.Y < r.Point.Y+r.H
}

func (r *Rect) Center() Point {
	return Point{X: r.Point.X + r.W/2, Y: r.Point.Y + r.H/2}
}

func (r *Rect) Intersect(other Rect) bool {
	return r.Point.X < other.Point.X+other.W && r.Point.X+r.W > other.Point.X && r.Point.Y < other.Point.Y+other.H && r.Point.Y+r.H > other.Point.Y
}

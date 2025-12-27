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
	return p.X > r.Point.X && p.X < r.Point.X+r.W-1 && p.Y > r.Point.Y && p.Y < r.Point.Y+r.H-1
}

func (r *Rect) Center() Point {
	return Point{X: r.Point.X + r.W/2, Y: r.Point.Y + r.H/2}
}

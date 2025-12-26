package domain

type Point struct {
	X int
	Y int
}

func (p *Point) Add(d Point) Point {
	return Point{X: p.X + d.X, Y: p.Y + d.Y}
}

func (p *Point) Sub(d Point) Point {
	return Point{X: p.X - d.X, Y: p.Y - d.Y}
}

var (
	DirUp        = Point{+0, -1}
	DirUpLeft    = Point{-1, -1}
	DirUpRight   = Point{+1, -1}
	DirDown      = Point{+0, +1}
	DirDownLeft  = Point{-1, +1}
	DirDownRight = Point{+1, +1}
	DirLeft      = Point{-1, +0}
	DirRight     = Point{+1, +0}
)

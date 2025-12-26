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
	DirSU = Point{+0, -1}
	DirUL = Point{-1, -1}
	DirUR = Point{+1, -1}
	DirSD = Point{+0, +1}
	DirDL = Point{-1, +1}
	DirDR = Point{+1, +1}
	DirLS = Point{-1, +0}
	DirRS = Point{+1, +0}
)

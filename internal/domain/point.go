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

var Dirs4 = []Point{DirSU, DirSD, DirLS, DirRS}
var Dirs8 = []Point{DirSU, DirSD, DirLS, DirRS, DirUL, DirUR, DirDL, DirDR}

func Manhattan(a, b Point) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func BresenhamLine(a, b Point) []Point {
	x0, y0 := a.X, a.Y
	x1, y1 := b.X, b.Y
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	var out []Point
	for {
		out = append(out, Point{X: x0, Y: y0})
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
	return out
}

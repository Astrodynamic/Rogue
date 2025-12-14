package game

// Direction represents a movement delta in grid coordinates.
// It is intentionally separate from Point to make "delta vs position" explicit in code.
type Direction struct {
	DX int
	DY int
}

var (
	DirUp    = Direction{DX: 0, DY: -1}
	DirDown  = Direction{DX: 0, DY: 1}
	DirLeft  = Direction{DX: -1, DY: 0}
	DirRight = Direction{DX: 1, DY: 0}
)

var (
	Dirs4 = [...]Direction{DirUp, DirDown, DirLeft, DirRight}
	Dirs8 = [...]Direction{
		{DX: 0, DY: -1},
		{DX: 0, DY: 1},
		{DX: -1, DY: 0},
		{DX: 1, DY: 0},
		{DX: -1, DY: -1},
		{DX: 1, DY: -1},
		{DX: -1, DY: 1},
		{DX: 1, DY: 1},
	}
)

func (p Point) Add(d Direction) Point { return Point{X: p.X + d.DX, Y: p.Y + d.DY} }

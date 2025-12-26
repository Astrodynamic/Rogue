package domain

type Corridor struct {
	Points []Point
}

func NewCorridor() Corridor {
	return Corridor{
		Points: make([]Point, 0),
	}
}

func (c *Corridor) AddPoint(p Point) {
	c.Points = append(c.Points, p)
}

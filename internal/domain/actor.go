package domain

type Actor struct {
	Point
	Stats
	Name string
}

func (a *Actor) Move(dir Point) {
	a.Point = a.Point.Add(dir)
}

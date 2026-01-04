package service

import "rogue/internal/domain"

type FOV struct {
	Radius int
}

func NewFOV() *FOV {
	return &FOV{
		Radius: domain.Radius,
	}
}

func (f *FOV) Compute(level *domain.Level, origin domain.Point) {
	if !level.Contains(origin) {
		return
	}

	level.ClearVisible()

	if f.computeRoomVisibility(level, origin) {
		return
	}

	f.computeLocalVisibility(level, origin)
	f.computeEntranceVisibility(level, origin)
}

func (f *FOV) Update(level *domain.Level, origin domain.Point) {
	f.Compute(level, origin)
	level.UpdateExplored()
}

func (f *FOV) computeRoomVisibility(level *domain.Level, origin domain.Point) bool {
	for _, room := range level.Rooms {
		if room.Contains(origin) {
			for y := room.Y; y < room.Y+room.H; y++ {
				for x := room.X; x < room.X+room.W; x++ {
					level.SetVisible(domain.Point{X: x, Y: y}, true)
				}
			}
			return true
		}
	}
	return false
}

func (f *FOV) computeLocalVisibility(level *domain.Level, origin domain.Point) {
	rad := f.Radius
	for dy := -rad; dy <= rad; dy++ {
		for dx := -rad; dx <= rad; dx++ {
			p := domain.Point{X: origin.X + dx, Y: origin.Y + dy}

			if domain.Manhattan(origin, p) > rad {
				continue
			}

			if f.hasLOS(level, origin, p) {
				level.SetVisible(p, true)
			}
		}
	}
}

func (f *FOV) computeEntranceVisibility(level *domain.Level, origin domain.Point) {
	for _, d := range domain.Dirs4 {
		neighbor := origin.Add(d)

		for _, room := range level.Rooms {
			if room.Contains(neighbor) {
				for x := room.X; x < room.X+room.W; x++ {
					f.castRay(level, origin, domain.Point{X: x, Y: room.Y})
					f.castRay(level, origin, domain.Point{X: x, Y: room.Y + room.H - 1})
				}

				for y := room.Y; y < room.Y+room.H; y++ {
					f.castRay(level, origin, domain.Point{X: room.X, Y: y})
					f.castRay(level, origin, domain.Point{X: room.X + room.W - 1, Y: y})
				}
				break
			}
		}
	}
}

func (f *FOV) castRay(level *domain.Level, from, to domain.Point) {
	line := domain.BresenhamLine(from, to)
	for _, p := range line {
		level.SetVisible(p, true)
		if p != from && level.BlocksSight(p) {
			return
		}
	}
}

func (f *FOV) hasLOS(level *domain.Level, from, to domain.Point) bool {
	line := domain.BresenhamLine(from, to)
	for _, p := range line {
		if p != from && level.BlocksSight(p) {
			return p == to
		}
	}
	return true
}

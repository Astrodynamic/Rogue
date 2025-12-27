package service

import (
	"rogue/internal/domain"
)

func (g *Generator) GenerateLevel(level *domain.Level) {
	level.Clear()
	g.GenerateRooms(level)
	g.GenerateCorridors(level)
	g.GenerateExit(level)
}

func (g *Generator) GenerateRooms(level *domain.Level) {
	level.Rooms = level.Rooms[:0]

	secW := level.W / 3
	secH := level.H / 3

	for sy := 0; sy < 3; sy++ {
		for sx := 0; sx < 3; sx++ {
			minW, minH := 6, 4
			maxW := secW - 2
			maxH := secH - 2
			if maxW < minW {
				maxW = minW
			}
			if maxH < minH {
				maxH = minH
			}
			w := g.rng.IntN(maxW-minW+1) + minW
			h := g.rng.IntN(maxH-minH+1) + minH

			secX := sx * secW
			secY := sy * secH
			maxX := secX + secW - w - 1
			maxY := secY + secH - h - 1
			x := secX + 1
			y := secY + 1
			if maxX > x {
				x = g.rng.IntN(maxX-x+1) + x
			}
			if maxY > y {
				y = g.rng.IntN(maxY-y+1) + y
			}

			level.AddRoom(domain.NewRoom(x, y, w, h))
		}
	}
}

func (g *Generator) GenerateCorridors(level *domain.Level) {
	level.Corridors = level.Corridors[:0]

	stack := []int{g.rng.IntN(len(level.Rooms))}
	visited := make([]bool, len(level.Rooms))
	visited[stack[0]] = true

	for len(stack) > 0 {
		cur := stack[len(stack)-1]

		neighbors := g.neighbors(cur)
		g.rng.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })

		next := -1
		for _, nb := range neighbors {
			if !visited[nb] {
				next = nb
				break
			}
		}

		if next == -1 {
			stack = stack[:len(stack)-1]
			continue
		}

		visited[next] = true
		stack = append(stack, next)

		g.GenerateCorridor(level, cur, next)
	}
}

func (g *Generator) neighbors(index int) []int {
	neighbors := make([]int, 0, 4)
	cx, cy := index%3, index/3
	if cx > 0 {
		neighbors = append(neighbors, index-1)
	}
	if cx < 2 {
		neighbors = append(neighbors, index+1)
	}
	if cy > 0 {
		neighbors = append(neighbors, index-3)
	}
	if cy < 2 {
		neighbors = append(neighbors, index+3)
	}
	return neighbors
}

func (g *Generator) GenerateCorridor(level *domain.Level, a, b int) {
	corridor := domain.NewCorridor()

	aRoom := level.Rooms[a]
	bRoom := level.Rooms[b]

	aPoint := aRoom.Center()
	bPoint := bRoom.Center()
	mPoint := aRoom.Center()

	if g.rng.IntN(2) == 0 {
		mPoint.X = bPoint.X
	} else {
		mPoint.Y = bPoint.Y
	}

	g.GenerateLine(&corridor, aPoint, mPoint)
	g.GenerateLine(&corridor, mPoint, bPoint)

	l, r := 0, len(corridor.Points)-1
	for l < r {
		if !aRoom.Contains(corridor.Points[l]) {
			break
		}
		l++
	}

	for l < r {
		if !bRoom.Contains(corridor.Points[r]) {
			break
		}
		r--
	}

	corridor.Points = corridor.Points[l : r+1]

	level.AddCorridor(corridor)
}

func (g *Generator) GenerateLine(corridor *domain.Corridor, a, b domain.Point) {
	x, y := a.X, a.Y

	dx := 0
	if b.X > x {
		dx = 1
	} else if b.X < x {
		dx = -1
	}

	dy := 0
	if b.Y > y {
		dy = 1
	} else if b.Y < y {
		dy = -1
	}

	for x != b.X || y != b.Y {
		corridor.AddPoint(domain.Point{X: x, Y: y})

		if x != b.X {
			x += dx
		} else if y != b.Y {
			y += dy
		}
	}
	corridor.AddPoint(b)
}

func (g *Generator) GenerateExit(level *domain.Level) {
	room := level.Rooms[g.rng.IntN(len(level.Rooms))]
	exit := domain.Point{
		X: g.rng.IntN(room.W-2) + room.X + 1,
		Y: g.rng.IntN(room.H-2) + room.Y + 1,
	}
	level.AddExit(exit)
}

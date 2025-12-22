package domain

// BresenhamLine returns points on a line from a to b inclusive.
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

	out := make([]Point, 0, dx+dy+1)
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

func (s *GameSession) ComputeVisible() [][]bool {
	w := len(s.Level.Tiles[0])
	h := len(s.Level.Tiles)
	vis := make2D(w, h, false)
	p := s.PlayerPos
	if !inBounds(p, w, h) {
		return vis
	}

	roomID, inRoom := s.roomAt(p)
	if inRoom {
		r := s.Level.Rooms[roomID].Bounds
		for y := r.Y; y < r.Y+r.H; y++ {
			for x := r.X; x < r.X+r.W; x++ {
				vis[y][x] = true
			}
		}
		return vis
	}

	// Corridor visibility: local LOS radius.
	const rad = 7
	for dy := -rad; dy <= rad; dy++ {
		for dx := -rad; dx <= rad; dx++ {
			t := Point{X: p.X + dx, Y: p.Y + dy}
			if !inBounds(t, w, h) {
				continue
			}
			if Manhattan(p, t) > rad {
				continue
			}
			if s.hasLOS(p, t) {
				vis[t.Y][t.X] = true
			}
		}
	}

	// If near a room entrance, reveal FOV into that room using ray casting + Bresenham.
	for _, d := range Dirs4 {
		np := p.Add(d)
		rid, ok := s.roomAt(np)
		if !ok {
			continue
		}
		r := s.Level.Rooms[rid].Bounds
		// Cast rays to room perimeter points (approx ray casting).
		for x := r.X; x < r.X+r.W; x++ {
			s.castRay(vis, p, Point{X: x, Y: r.Y})
			s.castRay(vis, p, Point{X: x, Y: r.Y + r.H - 1})
		}
		for y := r.Y; y < r.Y+r.H; y++ {
			s.castRay(vis, p, Point{X: r.X, Y: y})
			s.castRay(vis, p, Point{X: r.X + r.W - 1, Y: y})
		}
	}

	return vis
}

func (s *GameSession) castRay(vis [][]bool, from, to Point) {
	w := len(s.Level.Tiles[0])
	h := len(s.Level.Tiles)
	for _, pt := range BresenhamLine(from, to) {
		if !inBounds(pt, w, h) {
			return
		}
		vis[pt.Y][pt.X] = true
		if pt != from && s.Level.Tiles[pt.Y][pt.X].BlocksSight() {
			return
		}
	}
}

func (s *GameSession) hasLOS(from, to Point) bool {
	w := len(s.Level.Tiles[0])
	h := len(s.Level.Tiles)
	for _, pt := range BresenhamLine(from, to) {
		if !inBounds(pt, w, h) {
			return false
		}
		if pt != from && s.Level.Tiles[pt.Y][pt.X].BlocksSight() {
			return pt == to
		}
	}
	return true
}

func (s *GameSession) roomAt(p Point) (int, bool) {
	for i := range s.Level.Rooms {
		if s.Level.Rooms[i].Bounds.Contains(p) {
			return i, true
		}
	}
	return -1, false
}

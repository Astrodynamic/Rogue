package game

// ShortestPathNext returns the next step from start towards goal using 4-neighbor BFS.
// If no path exists (or already at goal), ok=false.
func ShortestPathNext(tiles [][]Tile, occupied map[Point]bool, start, goal Point) (next Point, ok bool) {
	if start == goal {
		return Point{}, false
	}
	w := len(tiles[0])
	h := len(tiles)
	if !inBounds(start, w, h) || !inBounds(goal, w, h) {
		return Point{}, false
	}

	prev := make(map[Point]Point, 128)
	q := make([]Point, 0, 128)
	seen := make(map[Point]bool, 128)

	q = append(q, start)
	seen[start] = true

	for qi := 0; qi < len(q); qi++ {
		cur := q[qi]
		for _, d := range Dirs4 {
			np := cur.Add(d)
			if !inBounds(np, w, h) || seen[np] {
				continue
			}
			if tiles[np.Y][np.X].BlocksMovement() {
				continue
			}
			if occupied[np] && np != goal {
				continue
			}
			seen[np] = true
			prev[np] = cur
			if np == goal {
				// Reconstruct 1 step.
				step := np
				for {
					p, has := prev[step]
					if !has {
						return Point{}, false
					}
					if p == start {
						return step, true
					}
					step = p
				}
			}
			q = append(q, np)
		}
	}
	return Point{}, false
}

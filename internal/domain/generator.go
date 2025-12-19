package domain

import "math/rand"

type Generator struct {
	rng *rand.Rand
}

func NewGenerator(rng *rand.Rand) *Generator { return &Generator{rng: rng} }

func (g *Generator) Generate(depth int, cfg Config) (Level, int, int) {
	cfg = cfg.Normalize()
	tiles := make2D(cfg.Width, cfg.Height, TileWall)
	rooms := make([]Room, 0, 9)

	secW := cfg.Width / 3
	secH := cfg.Height / 3

	roomID := 0
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
			w := g.rng.Intn(maxW-minW+1) + minW
			h := g.rng.Intn(maxH-minH+1) + minH

			secX := sx * secW
			secY := sy * secH
			maxX := secX + secW - w - 1
			maxY := secY + secH - h - 1
			x := secX + 1
			y := secY + 1
			if maxX > x {
				x = g.rng.Intn(maxX-x+1) + x
			}
			if maxY > y {
				y = g.rng.Intn(maxY-y+1) + y
			}

			r := Rect{X: x, Y: y, W: w, H: h}
			rooms = append(rooms, Room{ID: roomID, Bounds: r, FloorLoot: map[Point]Item{}})
			roomID++

			carveRoom(tiles, r)
		}
	}

	startIdx := g.rng.Intn(len(rooms))
	exitIdx := g.rng.Intn(len(rooms) - 1)
	if exitIdx >= startIdx {
		exitIdx++
	}
	rooms[startIdx].IsStart = true
	rooms[exitIdx].IsExit = true

	corridors := connectRooms(tiles, g.rng, rooms)

	exitPos := randomFloorInRoom(g.rng, tiles, rooms[exitIdx].Bounds)
	tiles[exitPos.Y][exitPos.X] = TileExit

	lvl := Level{
		Depth:     depth,
		Rooms:     rooms,
		Corridors: corridors,
		Tiles:     tiles,
	}
	return lvl, startIdx, exitIdx
}

func carveRoom(tiles [][]Tile, r Rect) {
	for y := r.Y; y < r.Y+r.H; y++ {
		for x := r.X; x < r.X+r.W; x++ {
			if x == r.X || y == r.Y || x == r.X+r.W-1 || y == r.Y+r.H-1 {
				tiles[y][x] = TileWall
			} else {
				tiles[y][x] = TileFloor
			}
		}
	}
}

func connectRooms(tiles [][]Tile, rng *rand.Rand, rooms []Room) []Corridor {
	// Randomized spanning tree on room centers ensures connectivity.
	n := len(rooms)
	visited := make([]bool, n)
	stack := []int{rng.Intn(n)}
	visited[stack[0]] = true

	corridors := make([]Corridor, 0, n-1)

	for len(stack) > 0 {
		cur := stack[len(stack)-1]

		neighbors := make([]int, 0, 4)
		cx, cy := cur%3, cur/3
		if cx > 0 {
			neighbors = append(neighbors, cur-1)
		}
		if cx < 2 {
			neighbors = append(neighbors, cur+1)
		}
		if cy > 0 {
			neighbors = append(neighbors, cur-3)
		}
		if cy < 2 {
			neighbors = append(neighbors, cur+3)
		}
		rng.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })

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

		a := rooms[cur].Bounds.Center()
		b := rooms[next].Bounds.Center()
		corridors = append(corridors, carveCorridor(tiles, rng, a, b))
	}

	return corridors
}

func carveCorridor(tiles [][]Tile, rng *rand.Rand, a, b Point) Corridor {
	// L-shaped corridor (randomly horizontal-first or vertical-first).
	tilesPath := make([]Point, 0, 64)
	hFirst := rng.Intn(2) == 0
	if hFirst {
		tilesPath = append(tilesPath, carveLine(tiles, a, Point{X: b.X, Y: a.Y})...)
		tilesPath = append(tilesPath, carveLine(tiles, Point{X: b.X, Y: a.Y}, b)...)
	} else {
		tilesPath = append(tilesPath, carveLine(tiles, a, Point{X: a.X, Y: b.Y})...)
		tilesPath = append(tilesPath, carveLine(tiles, Point{X: a.X, Y: b.Y}, b)...)
	}
	return Corridor{Tiles: tilesPath}
}

func carveLine(tiles [][]Tile, a, b Point) []Point {
	out := make([]Point, 0, 64)
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

	w := len(tiles[0])
	h := len(tiles)
	for x != b.X || y != b.Y {
		if inBounds(Point{X: x, Y: y}, w, h) && tiles[y][x] == TileWall {
			tiles[y][x] = TileCorridor
		} else if inBounds(Point{X: x, Y: y}, w, h) && tiles[y][x] == TileFloor {
			// keep room floor
		}
		out = append(out, Point{X: x, Y: y})

		if x != b.X {
			x += dx
		} else if y != b.Y {
			y += dy
		}
	}
	if inBounds(b, w, h) && tiles[b.Y][b.X] == TileWall {
		tiles[b.Y][b.X] = TileCorridor
	}
	out = append(out, b)
	return out
}

func randomFloorInRoom(rng *rand.Rand, tiles [][]Tile, r Rect) Point {
	// Guaranteed to find: room interior exists by construction.
	for {
		x := rng.Intn(r.W-2) + (r.X + 1)
		y := rng.Intn(r.H-2) + (r.Y + 1)
		if tiles[y][x] == TileFloor {
			return Point{X: x, Y: y}
		}
	}
}

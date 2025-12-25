package domain

type World struct {
	Level  *Level
	Player *Player
}

func NewWorld(h, w int) *World {
	lvl := NewLevel(w, h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if y == 0 || y == h-1 || x == 0 || x == w-1 {
				lvl.Tiles[y][x] = Tile{Kind: TileWall}
			} else {
				lvl.Tiles[y][x] = Tile{Kind: TileFloor}
			}
		}
	}
	player := &Player{Actor: Actor{Point: Point{X: w / 2, Y: h / 2}}}
	return &World{Level: lvl, Player: player}
}

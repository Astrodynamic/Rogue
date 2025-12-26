package domain

type Level struct {
	Width     int
	Height    int
	Tiles     [][]Tile
	Rooms     []Room
	Corridors []Corridor
}

func NewLevel(width, height int) *Level {
	tiles := make([][]Tile, height)
	for i := range tiles {
		tiles[i] = make([]Tile, width)
	}

	return &Level{
		Width: width, Height: height, Tiles: tiles,
		Rooms:     make([]Room, 0),
		Corridors: make([]Corridor, 0),
	}
}

func (l *Level) AddRoom(room Room) {
	l.Rooms = append(l.Rooms, room)
	l.FillRoom(room)
}

func (l *Level) FillRoom(room Room) {
	for y := room.Y; y < room.Y+room.H; y++ {
		for x := room.X; x < room.X+room.W; x++ {
			if x == room.X || y == room.Y || x == room.X+room.W-1 || y == room.Y+room.H-1 {
				l.Tiles[y][x] = Tile{Kind: TileWall}
			} else {
				l.Tiles[y][x] = Tile{Kind: TileFloor}
			}
		}
	}
}

func (l *Level) AddCorridor(corridor Corridor) {
	l.Corridors = append(l.Corridors, corridor)
	l.FillCorridor(corridor)
}

func (l *Level) FillCorridor(corridor Corridor) {
	for _, point := range corridor.Points {
		l.Tiles[point.Y][point.X] = Tile{Kind: TileCorridor}
	}

	for _, point := range corridor.Points {
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				if x == 0 && y == 0 {
					continue
				}

				if point.Y+y < 0 || point.Y+y >= l.Height || point.X+x < 0 || point.X+x >= l.Width {
					continue
				}

				tile := &l.Tiles[point.Y+y][point.X+x]
				if tile.Kind == TileNone {
					*tile = Tile{Kind: TileWall}
				}
			}
		}
	}
}

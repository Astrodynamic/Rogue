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
}

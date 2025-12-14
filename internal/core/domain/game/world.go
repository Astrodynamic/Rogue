package game

type Room struct {
	ID        int
	Bounds    Rect
	IsStart   bool
	IsExit    bool
	FloorLoot map[Point]Item
}

type Corridor struct {
	Tiles []Point
}

type Level struct {
	Depth     int
	Rooms     []Room
	Corridors []Corridor

	Tiles [][]Tile
}

type GameSession struct {
	Seed       int64
	LevelDepth int
	Level      Level

	Player    Character
	Backpack  Backpack
	PlayerPos Point

	NextItemID int64

	Enemies []Enemy

	// Fog of war.
	Explored [][]bool

	// Temporary effects, stored as deltas so expiry can revert.
	Effects []Effect

	Stats Stats

	// Recent messages for UI.
	Messages []string
}

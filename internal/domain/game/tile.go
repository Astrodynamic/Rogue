package game

type Tile uint8

const (
	TileVoid Tile = iota
	TileWall
	TileFloor
	TileCorridor
	TileExit
)

func (t Tile) BlocksMovement() bool {
	switch t {
	case TileWall, TileVoid:
		return true
	default:
		return false
	}
}

func (t Tile) BlocksSight() bool {
	switch t {
	case TileWall, TileVoid:
		return true
	default:
		return false
	}
}

package domain

type TileKind int

const (
	TileNone TileKind = iota
	TileWall
	TileFloor
	TileCorridor
)

type TileFlags uint8

const (
	TileVisible TileFlags = 1 << iota
)

type Tile struct {
	Kind  TileKind
	Flags TileFlags
}

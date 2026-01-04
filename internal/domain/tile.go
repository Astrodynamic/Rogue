package domain

type TileKind int

const (
	TileNone TileKind = iota
	TileWall
	TileFloor
	TileCorridor
	TileExit
)

type TileFlags uint8

const (
	TileVisible TileFlags = 1 << iota
	TileExplored
	TileOpaque
)

type Tile struct {
	Kind  TileKind
	Flags TileFlags
}

func (t *Tile) Clear() {
	t.SetKind(TileNone)
}

func (t *Tile) SetKind(kind TileKind) {
	t.Flags = 0
	switch kind {
	case TileNone:
		t.SetFlags(TileOpaque, true)
	case TileWall:
		t.SetFlags(TileOpaque, true)
	}

	t.Kind = kind
}

func (t *Tile) TestFlags(flags TileFlags) bool {
	return t.Flags&flags != 0
}

func (t *Tile) SetFlags(flags TileFlags, set bool) {
	if set {
		t.Flags |= flags
	} else {
		t.Flags &^= flags
	}
}

func (t *Tile) FlipFlags(flags TileFlags) {
	t.Flags ^= flags
}

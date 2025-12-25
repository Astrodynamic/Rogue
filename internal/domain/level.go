package domain

type Level struct {
	Width  int
	Height int
	Tiles  [][]Tile
}

func NewLevel(width, height int) *Level {
	tiles := make([][]Tile, height)
	for i := range tiles {
		tiles[i] = make([]Tile, width)
	}
	return &Level{Width: width, Height: height, Tiles: tiles}
}

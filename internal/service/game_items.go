package service

import "rogue/internal/domain"

func (g *Game) dropItemAtPosition(item domain.Item, pos domain.Point) bool {
	if !g.World.Level.Contains(pos) {
		return false
	}

	tile := g.World.Level.Tiles[pos.Y][pos.X]
	if tile.Kind != domain.TileFloor && tile.Kind != domain.TileCorridor {
		return false
	}

	if g.World.Level.GetItem(pos) != nil {
		return false
	}

	g.World.Level.AddItem(pos, item)
	return true
}

func (g *Game) findAdjacentDropPoint() domain.Point {
	return g.World.Level.GetAdjacentDropPoint(g.World.Player.Point)
}

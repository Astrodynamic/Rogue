package service

import "rogue/internal/domain"

func (g *Game) onMove(actor *domain.Actor, dir domain.Point) {
	next := actor.Point.Add(dir)

	if !g.World.Level.Contains(next) {
		return
	}

	switch g.World.Level.Tiles[next.Y][next.X].Kind {
	case domain.TileWall:
		return
	case domain.TileExit:
		NewGenerator().Generate(g.World)
		return
	default:
		actor.Move(dir)
	}
}

package service

import "rogue/internal/domain"

func (g *Game) onMove(actor *domain.Actor, dir domain.Point) {
	next := actor.Point.Add(dir)

	if !g.World.Level.Contains(next) {
		return
	}

	if g.World.Level.Tiles[next.Y][next.X].Kind == domain.TileWall {
		return
	}

	if g.World.Level.Tiles[next.Y][next.X].Kind == domain.TileExit {
		NewGenerator().Generate(g.World)
		return
	}

	actor.Move(dir)
}

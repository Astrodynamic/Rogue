package service

import "rogue/internal/domain"

func (g *Game) onMove(actor *domain.Actor, dir domain.Point) {
	next := actor.Point.Add(dir)

	if next.X < 0 || next.X >= g.World.Level.Width || next.Y < 0 || next.Y >= g.World.Level.Height {
		return
	}

	if g.World.Level.Tiles[next.Y][next.X].Kind == domain.TileWall {
		return
	}

	actor.Move(dir)
}

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
	default:
		actor.Move(dir)
		g.pickupItem(actor, next)
	}

	actor.TickEffects()
	g.UpdateVisibility()
}

func (g *Game) pickupItem(actor *domain.Actor, pos domain.Point) {
	item := g.World.Level.GetItem(pos)
	if item == nil {
		return
	}

	if !actor.Backpack.HasSpace(item) {
		return
	}

	if actor.Backpack.Add(item) {
		g.World.Level.RemoveItem(pos)
	}
}

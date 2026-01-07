package service

import (
	"rogue/internal/domain"
)

func (g *Game) onMove(actor *domain.Actor, dir domain.Point) {
	next := actor.Point.Add(dir)

	if !g.World.Level.Contains(next) {
		return
	}

	enemy := g.World.Level.GetEnemy(next)
	if enemy != nil && enemy.IsAlive() {
		g.initiateCombat(actor, enemy, next)
		return
	}

	switch g.World.Level.Tiles[next.Y][next.X].Kind {
	case domain.TileWall:
		return
	case domain.TileExit:
		if g.World.GetDepth() >= domain.Depth {
			g.handleGameCompletion()
			return
		}
		g.SaveStatistics()
		g.SaveGameState()
		NewGenerator().Generate(g.World)
		g.UpdateVisibility()
		return
	default:
		actor.Move(dir)
		g.RecordTileTraveled()
		g.pickupItem(actor, next)
	}

	actor.TickEffects()
	g.UpdateVisibility()

	g.ProcessEnemyTurns()

	if g.World.Player.Health <= 0 {
		g.handlePlayerDeath()
	}
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
		if item.Type() == domain.ItemTreasure {
			if treasure, ok := item.(*domain.Treasure); ok {
				g.RecordTreasureCollected(treasure.Value)
			}
		}
		g.World.Level.RemoveItem(pos)
	}
}

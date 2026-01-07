package service

import (
	"fmt"
	"rogue/internal/domain"
)

func (g *Game) onMove(actor *domain.Actor, dir domain.Point) {

	if actor.State == domain.ActorStateSleep {
		g.ui.AddLog("Awake")
		actor.State = domain.ActorStateNormal
		actor.TickEffects()
		g.UpdateVisibility()
		g.ProcessEnemyTurns()
		g.checkPlayerHealth()
		return
	}

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
		g.ui.AddLog("Blocked")

		actor.TickEffects()
		g.UpdateVisibility()
		g.ProcessEnemyTurns()
		g.checkPlayerHealth()
		return
	case domain.TileExit:
		if g.World.GetDepth() >= domain.WorldConfig.Depth {
			g.handleGameCompletion()
			return
		}
		g.ui.AddLog(fmt.Sprintf("Level %d", g.World.GetDepth()+1))

		g.World.GameState.AdjustDifficulty(g.World.Player.Health, g.World.Player.MaxHealth)

		g.SaveStatistics()
		g.SaveGameState()
		g.World.GameState.AdvanceLevel()
		g.generator.Generate(g.World)
		g.UpdateVisibility()
		g.ui.AddLog(fmt.Sprintf("Enter %d", g.World.GetDepth()))
		return
	default:
		actor.Move(dir)
		g.RecordTileTraveled()
		g.pickupItem(actor, next)
	}

	actor.TickEffects()
	g.UpdateVisibility()

	g.ProcessEnemyTurns()
	g.checkPlayerHealth()
}

func (g *Game) pickupItem(actor *domain.Actor, pos domain.Point) {
	item := g.World.Level.GetItem(pos)
	if item == nil {
		return
	}

	if !actor.Backpack.HasSpace(item) {
		g.ui.AddLog("Backpack full")
		return
	}

	if !actor.Backpack.Add(item) {
		return
	}

	if item.Type() == domain.ItemTreasure {
		if treasure, ok := item.(*domain.Treasure); ok {
			g.RecordTreasureCollected(treasure.Value)
			g.ui.AddLog(fmt.Sprintf("Got $%d", treasure.Value))
		}
	} else {
		g.ui.AddLog("Got " + item.Name())
	}
	g.World.Level.RemoveItem(pos)
}

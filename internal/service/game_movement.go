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
		g.SaveStatistics()
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

func (g *Game) initiateCombat(playerActor *domain.Actor, enemy domain.Enemy, enemyPos domain.Point) {
	enemyActor := enemy.GetActor()

	resolver := NewGameCombatResolver(g)

	if vampire, ok := enemy.(*domain.Vampire); ok {
		if !vampire.FirstHitMissed {
			vampire.MarkFirstHitMissed()
		}
	}

	result := playerActor.AttackWithResolver(enemyActor, resolver)
	g.RecordHitDealt()

	if result.Hit {
		if result.Killed {
			g.handleEnemyDeath(enemy, enemyPos)
			g.RecordEnemyDefeated()
			return
		}
	}

	if enemy.IsAlive() && enemyActor.State != domain.ActorStateSleep {
		playerActor := &g.World.Player.Actor

		switch e := enemy.(type) {
		case *domain.Ogre:
			if e.Resting {
				enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
				g.RecordHitReceived()
				if enemyResult.Killed {
					g.handlePlayerDeath()
					return
				}
				e.StartRest()
				return
			}
		case *domain.SnakeMage:
			enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
			g.RecordHitReceived()
			if enemyResult.Hit {
				if g.rng.IntN(100) < domain.SnakeMageSleepChance {
					sleepEffect := domain.NewSleepEffect(1)
					playerActor.AddEffect(sleepEffect)
				}
			}
			if enemyResult.Killed {
				g.handlePlayerDeath()
				return
			}
		case *domain.Vampire:
			enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
			g.RecordHitReceived()
			if enemyResult.Hit {
				reductionEffect := domain.NewMaxHealthEffect(-domain.VampireMaxHealthReduction, -1)
				playerActor.AddEffect(reductionEffect)
			}
			if enemyResult.Killed {
				g.handlePlayerDeath()
				return
			}
		default:
			enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
			g.RecordHitReceived()
			if enemyResult.Killed {
				g.handlePlayerDeath()
				return
			}
		}

		if ogre, ok := enemy.(*domain.Ogre); ok {
			ogre.StartRest()
		}
	}

	g.UpdateVisibility()
}

func (g *Game) handleEnemyDeath(enemy domain.Enemy, enemyPos domain.Point) {
	depth := g.World.GetDepth()
	generator := NewGenerator()
	treasure := generator.GenerateTreasureFromEnemyWithDepth(enemy, depth)

	if treasure != nil {
		g.World.Level.AddItem(enemyPos, treasure)
	}

	g.World.Level.RemoveEnemy(enemyPos)
}

func (g *Game) handlePlayerDeath() {
	g.SaveStatistics()
	world := domain.NewWorld(domain.Width, domain.Height)
	g.World = world
	NewGenerator().Generate(g.World)
	g.UpdateVisibility()
}

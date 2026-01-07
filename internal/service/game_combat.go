package service

import (
	"fmt"
	"rogue/internal/domain"
)

func (g *Game) initiateCombat(playerActor *domain.Actor, enemy domain.Enemy, enemyPos domain.Point) {
	enemyActor := enemy.GetActor()
	enemyName := enemy.Name()
	resolver := NewGameCombatResolver(g)

	g.handleVampireFirstHit(enemy)

	result := playerActor.AttackWithResolver(enemyActor, resolver)
	g.RecordHitDealt()

	if result.Hit && result.Killed {
		g.ui.AddLog(fmt.Sprintf("Defeated %s!", enemyName))
		g.handleEnemyDeath(enemy, enemyPos)
		g.RecordEnemyDefeated()
		return
	}

	if result.Hit {
		g.ui.AddLog(fmt.Sprintf("Hit %s for %d damage", enemyName, result.Damage))
	} else {
		g.ui.AddLog(fmt.Sprintf("Missed %s", enemyName))
	}

	if enemy.IsAlive() && enemyActor.State != domain.ActorStateSleep {
		g.processEnemyAttack(enemy, enemyActor, resolver)
	}

	g.UpdateVisibility()
}

func (g *Game) handleVampireFirstHit(enemy domain.Enemy) {
	if vampire, ok := enemy.(*domain.Vampire); ok {
		if !vampire.FirstHitMissed {
			vampire.MarkFirstHitMissed()
			g.ui.AddLog("Vampire first hit missed!")
		}
	}
}

func (g *Game) processEnemyAttack(enemy domain.Enemy, enemyActor *domain.Actor, resolver domain.CombatResolver) {
	playerActor := &g.World.Player.Actor

	switch e := enemy.(type) {
	case *domain.Ogre:
		g.processOgreAttack(e, enemyActor, playerActor, resolver)
	case *domain.SnakeMage:
		g.processSnakeMageAttack(e, enemyActor, playerActor, resolver)
	case *domain.Vampire:
		g.processVampireAttack(e, enemyActor, playerActor, resolver)
	default:
		g.processDefaultEnemyAttack(enemyActor, playerActor, resolver)
	}
}

func (g *Game) processOgreAttack(ogre *domain.Ogre, enemyActor, playerActor *domain.Actor, resolver domain.CombatResolver) {
	if ogre.Resting {
		enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
		g.RecordHitReceived()
		if enemyResult.Hit {
			g.ui.AddLog(fmt.Sprintf("Ogre hit you for %d damage", enemyResult.Damage))
		} else {
			g.ui.AddLog("Ogre missed!")
		}
		if enemyResult.Killed {
			g.handlePlayerDeath()
			return
		}
		ogre.StartRest()
		return
	}

	enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
	g.RecordHitReceived()
	if enemyResult.Hit {
		g.ui.AddLog(fmt.Sprintf("Ogre hit you for %d damage", enemyResult.Damage))
	} else {
		g.ui.AddLog("Ogre missed!")
	}
	if enemyResult.Killed {
		g.handlePlayerDeath()
		return
	}
	ogre.StartRest()
}

func (g *Game) processSnakeMageAttack(snakeMage *domain.SnakeMage, enemyActor, playerActor *domain.Actor, resolver domain.CombatResolver) {
	enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
	g.RecordHitReceived()
	if enemyResult.Hit {
		g.ui.AddLog(fmt.Sprintf("Snake-Mage hit you for %d damage", enemyResult.Damage))
		if g.rng.IntN(100) < domain.SnakeMageSleepChance {
			sleepEffect := domain.NewSleepEffect(1)
			playerActor.AddEffect(sleepEffect)
			g.ui.AddLog("Put to sleep by Snake-Mage!")
		}
	} else {
		g.ui.AddLog("Snake-Mage missed!")
	}
	if enemyResult.Killed {
		g.handlePlayerDeath()
	}
}

func (g *Game) processVampireAttack(vampire *domain.Vampire, enemyActor, playerActor *domain.Actor, resolver domain.CombatResolver) {
	enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
	g.RecordHitReceived()
	if enemyResult.Hit {
		g.ui.AddLog(fmt.Sprintf("Vampire hit you for %d damage", enemyResult.Damage))
		reductionEffect := domain.NewMaxHealthEffect(-domain.VampireMaxHealthReduction, -1)
		playerActor.AddEffect(reductionEffect)
		g.ui.AddLog(fmt.Sprintf("Maximum health reduced by %d!", domain.VampireMaxHealthReduction))
	} else {
		g.ui.AddLog("Vampire missed!")
	}
	if enemyResult.Killed {
		g.handlePlayerDeath()
	}
}

func (g *Game) processDefaultEnemyAttack(enemyActor, playerActor *domain.Actor, resolver domain.CombatResolver) {
	enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
	g.RecordHitReceived()
	if enemyResult.Hit {
		g.ui.AddLog(fmt.Sprintf("Received %d damage", enemyResult.Damage))
	} else {
		g.ui.AddLog("Enemy missed!")
	}
	if enemyResult.Killed {
		g.handlePlayerDeath()
	}
}

func (g *Game) handleEnemyDeath(enemy domain.Enemy, enemyPos domain.Point) {
	depth := g.World.GetDepth()
	treasure := g.generator.GenerateTreasureFromEnemyWithDepth(enemy, depth)

	if treasure != nil {
		g.World.Level.AddItem(enemyPos, treasure)
		g.ui.AddLog(fmt.Sprintf("Dropped %s (value: %d)", treasure.Name(), treasure.Value))
	}

	g.World.Level.RemoveEnemy(enemyPos)
}

func (g *Game) handlePlayerDeath() {
	g.ui.AddLog("You died! Restarting game...")
	g.SaveStatistics()
	world := domain.NewWorld(domain.Width, domain.Height)
	g.World = world
	g.generator.Generate(g.World)
	g.UpdateVisibility()
}

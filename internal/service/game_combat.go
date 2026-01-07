package service

import (
	"fmt"
	"rogue/internal/domain"
)

func (g *Game) initiateCombat(playerActor *domain.Actor, enemy domain.Enemy, enemyPos domain.Point) {
	enemyActor := enemy.GetActor()
	enemyName := enemy.Name()
	resolver := NewGameCombatResolver(g)

	if mimic, ok := enemy.(*domain.Mimic); ok {
		if !mimic.IsRevealed() {
			mimic.Reveal()
		}
	}

	result := playerActor.AttackWithResolver(enemyActor, resolver)
	g.RecordHitDealt()

	if vampire, ok := enemy.(*domain.Vampire); ok {
		if !vampire.FirstHitMissed {
			vampire.MarkFirstHitMissed()
			if !result.Hit {
				g.ui.AddLog("Dodge!")
			}
		}
	}

	if result.Hit && result.Killed {
		g.ui.AddLog(fmt.Sprintf("Killed %s", enemyName))
		g.handleEnemyDeath(enemy, enemyPos)
		g.RecordEnemyDefeated()
		return
	}

	if result.Hit {
		g.ui.AddLog(fmt.Sprintf("Hit %s %ddmg", enemyName, result.Damage))
	} else {
		g.ui.AddLog("Miss")
	}

	if enemy.IsAlive() && enemyActor.State != domain.ActorStateSleep {
		g.processEnemyAttack(enemy, enemyActor, resolver)
	}

	g.UpdateVisibility()
}

func (g *Game) logAttackResult(result domain.CombatResult) {
	if result.Hit {
		g.ui.AddLog(fmt.Sprintf("Hit %ddmg", result.Damage))
	} else {
		g.ui.AddLog("Miss")
	}
}

func (g *Game) processEnemyAttack(enemy domain.Enemy, actor *domain.Actor, resolver domain.CombatResolver) {
	playerActor := &g.World.Player.Actor
	enemyResult := actor.AttackWithResolver(playerActor, resolver)

	g.RecordHitReceived()
	g.logAttackResult(enemyResult)

	g.applyEnemyAttackEffects(enemy, enemyResult, playerActor)

	if enemyResult.Killed {
		g.handlePlayerDeath()
		return
	}

	if ogre, ok := enemy.(*domain.Ogre); ok {
		ogre.StartRest()
	}
}

func (g *Game) applyEnemyAttackEffects(enemy domain.Enemy, result domain.CombatResult, playerActor *domain.Actor) {
	if !result.Hit {
		return
	}

	switch enemy.(type) {
	case *domain.SnakeMage:
		if g.rng.IntN(domain.Combat.PercentBase) < domain.SnakeMageSleepChance {
			playerActor.State = domain.ActorStateSleep
			g.ui.AddLog("Sleep!")
		}
	case *domain.Vampire:
		reductionEffect := domain.NewMaxHealthEffect(-domain.VampireMaxHealthReduction, -1)
		playerActor.AddEffect(reductionEffect)
		g.ui.AddLog(fmt.Sprintf("MaxHP-%d", domain.VampireMaxHealthReduction))
	}
}

func (g *Game) handleEnemyDeath(enemy domain.Enemy, enemyPos domain.Point) {
	depth := g.World.GetDepth()
	treasure := g.generator.GenerateTreasureFromEnemyWithDepth(enemy, depth)

	if treasure != nil {
		g.World.Level.AddItem(enemyPos, treasure)
		g.ui.AddLog(fmt.Sprintf("Drop $%d", treasure.Value))
	}

	g.World.Level.RemoveEnemy(enemyPos)
}

func (g *Game) handlePlayerDeath() {
	g.ui.AddLog("Died!")
	g.SaveStatistics()
	world := domain.NewWorld(domain.WorldConfig.Width, domain.WorldConfig.Height, g.World.Player.Name)
	g.World = world
	g.generator.Generate(g.World)
	g.UpdateVisibility()
}

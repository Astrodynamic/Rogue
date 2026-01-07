package service

import (
	"fmt"
	"rogue/internal/domain"
)

func (g *Game) initiateCombat(pActor *domain.Actor, enemy domain.Enemy, ePos domain.Point) {
	eActor := enemy.GetActor()
	enemyName := enemy.Name()
	resolver := NewCombatRes(g)

	if mimic, ok := enemy.(*domain.Mimic); ok {
		if !mimic.IsRevealed() {
			mimic.Reveal()
		}
	}

	result := pActor.AttackWithResolver(eActor, resolver)
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
		g.handleEnemyDeath(enemy, ePos)
		g.RecordEnemyDefeated()
		return
	}

	if result.Hit {
		g.ui.AddLog(fmt.Sprintf("Hit %s %ddmg", enemyName, result.Damage))
	} else {
		g.ui.AddLog("Miss")
	}

	if enemy.IsAlive() && eActor.State != domain.ActorStateSleep {
		g.processEnemyAttack(enemy, eActor, resolver)
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

func (g *Game) processEnemyAttack(enemy domain.Enemy, eActor *domain.Actor, resolver domain.CombatResolver) {
	pActor := &g.World.Player.Actor
	result := eActor.AttackWithResolver(pActor, resolver)

	g.RecordHitReceived()
	g.logAttackResult(result)

	g.applyEnemyAttackEffects(enemy, result, pActor)

	if result.Killed {
		g.handlePlayerDeath()
		return
	}

	if ogre, ok := enemy.(*domain.Ogre); ok {
		ogre.StartRest()
	}
}

func (g *Game) applyEnemyAttackEffects(enemy domain.Enemy, result domain.CombatResult, pActor *domain.Actor) {
	if !result.Hit {
		return
	}

	switch enemy.(type) {
	case *domain.SnakeMage:
		if g.rng.IntN(domain.Combat.PercentBase) < domain.SnakeMageSleepCh {
			pActor.State = domain.ActorStateSleep
			g.ui.AddLog("Sleep!")
		}
	case *domain.Vampire:
		reductionEffect := domain.NewMaxHPEffect(-domain.VampireMaxHPRed, -1)
		pActor.AddEffect(reductionEffect)
		g.ui.AddLog(fmt.Sprintf("MaxHP-%d", domain.VampireMaxHPRed))
	}
}

func (g *Game) handleEnemyDeath(enemy domain.Enemy, ePos domain.Point) {
	depth := g.World.GetDepth()
	treasure := g.gen.GenerateTreasureFromEnemyWithDepth(enemy, depth)

	if treasure != nil {
		g.World.Level.AddItem(ePos, treasure)
		g.ui.AddLog(fmt.Sprintf("Drop $%d", treasure.Value))
	}

	g.World.Level.RemoveEnemy(ePos)
}

func (g *Game) handlePlayerDeath() {
	g.ui.AddLog("Died!")
	g.SaveStatistics()
	world := domain.NewWorld(domain.WorldConfig.Width, domain.WorldConfig.Height, g.World.Player.Name)
	g.World = world
	g.gen.Generate(g.World)
	g.UpdateVisibility()
}

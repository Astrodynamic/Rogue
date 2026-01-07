package service

import (
	"rogue/internal/domain"
)

func (g *Game) initiateCombat(playerActor *domain.Actor, enemy domain.Enemy, enemyPos domain.Point) {
	enemyActor := enemy.GetActor()
	resolver := NewGameCombatResolver(g)

	g.handleVampireFirstHit(enemy)

	result := playerActor.AttackWithResolver(enemyActor, resolver)
	g.RecordHitDealt()

	if result.Hit && result.Killed {
		g.handleEnemyDeath(enemy, enemyPos)
		g.RecordEnemyDefeated()
		return
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
		if enemyResult.Killed {
			g.handlePlayerDeath()
			return
		}
		ogre.StartRest()
		return
	}

	enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
	g.RecordHitReceived()
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
		if g.rng.IntN(100) < domain.SnakeMageSleepChance {
			sleepEffect := domain.NewSleepEffect(1)
			playerActor.AddEffect(sleepEffect)
		}
	}
	if enemyResult.Killed {
		g.handlePlayerDeath()
	}
}

func (g *Game) processVampireAttack(vampire *domain.Vampire, enemyActor, playerActor *domain.Actor, resolver domain.CombatResolver) {
	enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
	g.RecordHitReceived()
	if enemyResult.Hit {
		reductionEffect := domain.NewMaxHealthEffect(-domain.VampireMaxHealthReduction, -1)
		playerActor.AddEffect(reductionEffect)
	}
	if enemyResult.Killed {
		g.handlePlayerDeath()
	}
}

func (g *Game) processDefaultEnemyAttack(enemyActor, playerActor *domain.Actor, resolver domain.CombatResolver) {
	enemyResult := enemyActor.AttackWithResolver(playerActor, resolver)
	g.RecordHitReceived()
	if enemyResult.Killed {
		g.handlePlayerDeath()
	}
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

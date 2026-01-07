package service

import (
	"rogue/internal/domain"
)

type GameCombatResolver struct {
	game *Game
}

func NewGameCombatResolver(game *Game) *GameCombatResolver {
	return &GameCombatResolver{game: game}
}

func (r *GameCombatResolver) GetWeaponStrength(actor *domain.Actor) int {
	playerActor := &r.game.World.Player.Actor
	if actor == playerActor {
		if r.game.World.Player.Equipment != nil {
			weapon := r.game.World.Player.Equipment.Get(domain.ActorPartHand)
			if weapon != nil {
				return weapon.GetStats().Strength
			}
		}
	}
	return 0
}

func (r *GameCombatResolver) GetRandomGenerator() domain.RandomGenerator {
	return r.game.rng
}

func (r *GameCombatResolver) GetAttackModifiers(attacker *domain.Actor, target *domain.Actor) domain.AttackModifiers {
	modifiers := domain.AttackModifiers{
		GuaranteedHit: false,
		FirstHitMiss:  false,
	}

	if target != nil {
		enemies := r.game.World.Level.GetAllEnemies()
		for _, enemyWithPos := range enemies {
			enemyActor := enemyWithPos.Enemy.GetActor()
			if enemyActor == target {
				if vampire, ok := enemyWithPos.Enemy.(*domain.Vampire); ok {
					modifiers.FirstHitMiss = !vampire.FirstHitMissed
				}
				break
			}
		}
	}

	if attacker != nil {
		enemies := r.game.World.Level.GetAllEnemies()
		for _, enemyWithPos := range enemies {
			enemyActor := enemyWithPos.Enemy.GetActor()
			if enemyActor == attacker {
				if ogre, ok := enemyWithPos.Enemy.(*domain.Ogre); ok {
					if ogre.Resting {
						modifiers.GuaranteedHit = true
					}
				}
				break
			}
		}
	}

	return modifiers
}

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
		if e := r.findEnemyByActor(target); e != nil {
			if v, ok := e.Enemy.(*domain.Vampire); ok {
				modifiers.FirstHitMiss = !v.FirstHitMissed
			}
		}
	}

	if attacker != nil {
		if e := r.findEnemyByActor(attacker); e != nil {
			if o, ok := e.Enemy.(*domain.Ogre); ok && !o.Resting && o.RestTurns == 0 {
				modifiers.GuaranteedHit = true
			}
		}
	}

	return modifiers
}

func (r *GameCombatResolver) findEnemyByActor(actor *domain.Actor) *domain.EnemyWithPos {
	enemies := r.game.World.Level.GetAllEnemies()
	for i := range enemies {
		if enemies[i].Enemy.GetActor() == actor {
			return &enemies[i]
		}
	}
	return nil
}

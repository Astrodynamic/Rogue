package service

import (
	"rogue/internal/domain"
)

type CombatRes struct {
	game *Game
}

func NewCombatRes(game *Game) *CombatRes {
	return &CombatRes{game: game}
}

func (r *CombatRes) GetWeaponStrength(actor *domain.Actor) int {
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

func (r *CombatRes) GetRandomGenerator() domain.RandomGenerator {
	return r.game.rng
}

func (r *CombatRes) GetAttackModifiers(attacker *domain.Actor, target *domain.Actor) domain.AttackModifiers {
	modifiers := domain.AttackModifiers{
		GuaranteedHit: false,
		FirstHitMiss:  false,
	}

	if target != nil {
		if e := r.findEnemyByActor(target); e != nil {
			if v, ok := e.(*domain.Vampire); ok {
				modifiers.FirstHitMiss = !v.FirstHitMissed
			}
		}
	}

	if attacker != nil {
		if e := r.findEnemyByActor(attacker); e != nil {
			if o, ok := e.(*domain.Ogre); ok && !o.Resting && o.RestTurns == 0 {
				modifiers.GuaranteedHit = true
			}
		}
	}

	return modifiers
}

func (r *CombatRes) findEnemyByActor(actor *domain.Actor) domain.Enemy {
	enemies := r.game.World.Level.GetAllEnemies()
	for _, enemy := range enemies {
		if enemy.GetActor() == actor {
			return enemy
		}
	}
	return nil
}

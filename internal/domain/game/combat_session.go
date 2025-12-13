package game

import "math/rand"

func (s *GameSession) playerAttack(rng *rand.Rand, enemyIdx int) {
	e := &s.Enemies[enemyIdx]
	if e.Type == EnemyVampire && e.VampireFirstHitAlwaysMiss {
		e.VampireFirstHitAlwaysMiss = false
		s.addMsg("Your first strike against the vampire misses!")
		return
	}
	if !RollHit(rng, s.Player.Dexterity, e.Dexterity) {
		s.addMsg("You miss the " + e.Type.Name() + ".")
		return
	}
	s.Stats.HitsDealt++
	weaponBonus := 0
	if wid := s.Player.Equipment.Get(SlotWeapon); wid != 0 {
		if w, ok := s.Backpack.FindByID(wid); ok && w.Type == ItemWeapon {
			weaponBonus = w.Strength
		}
	}
	dmg := Damage(rng, s.Player.Strength, weaponBonus)
	e.Health -= dmg
	if e.Health > 0 {
		s.addMsg("You hit " + e.Type.Name() + " for " + itoa(dmg) + " (hp " + itoa(e.Health) + ").")
	} else {
		s.addMsg("You hit " + e.Type.Name() + " for " + itoa(dmg) + ".")
	}
	if e.Health <= 0 {
		drop := TreasureDrop(rng, *e)
		s.Backpack.AddTreasure(drop)
		s.Stats.TotalTreasure += drop
		s.Stats.DefeatedEnemies++
		xp := XPDrop(*e)
		gained := s.Player.GainXP(xp)
		s.addMsg(e.Type.Name() + " defeated. Treasure +" + itoa(drop) + ", XP +" + itoa(xp) + ".")
		if gained > 0 {
			s.addMsg("Level up! L" + itoa(s.Player.Level) + "  HP " + itoa(s.Player.Health) + "/" + itoa(s.Player.MaxHealth) + "  Dex " + itoa(s.Player.Dexterity) + "  Str " + itoa(s.Player.Strength) + ".")
		}
		// remove enemy
		copy(s.Enemies[enemyIdx:], s.Enemies[enemyIdx+1:])
		s.Enemies = s.Enemies[:len(s.Enemies)-1]
	}
}

func (s *GameSession) enemyAttack(rng *rand.Rand, e *Enemy) {
	guaranteed := false
	if e.Type == EnemyOgre && e.OgreGuaranteedHitNext {
		guaranteed = true
		e.OgreGuaranteedHitNext = false
	}

	hit := guaranteed || RollHit(rng, e.Dexterity, s.Player.Dexterity)
	if !hit {
		s.addMsg(e.Type.Name() + " misses.")
		return
	}
	s.Stats.HitsReceived++
	dmg := Damage(rng, e.Strength, 0)
	s.Player.ApplyDamage(dmg)
	s.addMsg(e.Type.Name() + " hits you for " + itoa(dmg) + " (hp " + itoa(s.Player.Health) + ").")

	switch e.Type {
	case EnemyVampire:
		// Reduce max health on successful attack.
		if s.Player.MaxHealth > 1 {
			s.Player.MaxHealth--
			if s.Player.Health > s.Player.MaxHealth {
				s.Player.Health = s.Player.MaxHealth
			}
			s.addMsg("Vampire drains your max HP!")
		}
	case EnemySnakeMage:
		// Chance to put player to sleep for one turn.
		if rng.Float64() < 0.30 {
			s.addEffect(Effect{Type: EffectSleep, TurnsLeft: 1, Delta: 0})
			s.addMsg("You fall asleep!")
		}
	case EnemyOgre:
		e.OgreRestTurns = 1
		e.OgreGuaranteedHitNext = true
	}

	if s.Player.IsDead() {
		s.addMsg("You died.")
	}
}

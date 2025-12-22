package domain

import "math/rand"

func (s *GameSession) useFood(idx int) {
	it, ok := s.Backpack.RemoveAt(ItemFood, idx)
	if !ok {
		s.addMsg("No such food.")
		return
	}
	s.Player.Heal(it.Health)
	s.Stats.FoodConsumed++
	s.addMsg("You eat food.")
}

func (s *GameSession) useElixir(rng *rand.Rand, idx int) {
	it, ok := s.Backpack.RemoveAt(ItemElixir, idx)
	if !ok {
		s.addMsg("No such elixir.")
		return
	}
	s.Stats.ElixirsDrunk++
	switch it.Subtype {
	case SubtypeDexterity:
		s.Player.Dexterity += it.Dexterity
		s.addEffect(Effect{Type: EffectDexBuff, TurnsLeft: max(1, it.Duration), Delta: it.Dexterity})
	case SubtypeStrength:
		s.Player.Strength += it.Strength
		s.addEffect(Effect{Type: EffectStrBuff, TurnsLeft: max(1, it.Duration), Delta: it.Strength})
	case SubtypeMaxHealth:
		s.Player.IncreaseMaxHealth(it.MaxHealth)
		s.addEffect(Effect{Type: EffectMaxHPBuff, TurnsLeft: max(1, it.Duration), Delta: it.MaxHealth})
	default:
		// fallback random buff
		switch rng.Intn(3) {
		case 0:
			s.Player.Dexterity += 1
			s.addEffect(Effect{Type: EffectDexBuff, TurnsLeft: 12, Delta: 1})
		case 1:
			s.Player.Strength += 1
			s.addEffect(Effect{Type: EffectStrBuff, TurnsLeft: 12, Delta: 1})
		default:
			s.Player.IncreaseMaxHealth(1)
			s.addEffect(Effect{Type: EffectMaxHPBuff, TurnsLeft: 12, Delta: 1})
		}
	}
	s.addMsg("You drink an elixir.")
}

func (s *GameSession) useScroll(idx int) {
	it, ok := s.Backpack.RemoveAt(ItemScroll, idx)
	if !ok {
		s.addMsg("No such scroll.")
		return
	}
	s.Stats.ScrollsRead++
	switch it.Subtype {
	case SubtypeDexterity:
		s.Player.Dexterity += it.Dexterity
	case SubtypeStrength:
		s.Player.Strength += it.Strength
	case SubtypeMaxHealth:
		s.Player.IncreaseMaxHealth(it.MaxHealth)
	default:
		s.Player.Dexterity += 1
	}
	s.addMsg("You read a scroll.")
}

func (s *GameSession) equipWeapon(rng *rand.Rand, idx int) {
	weapons := s.Backpack.List(ItemWeapon)
	if idx < 0 || idx >= len(weapons) {
		s.addMsg("No such weapon.")
		return
	}
	chosen := weapons[idx]
	if chosen.ID != 0 && s.Player.Equipment.Get(SlotWeapon) == chosen.ID {
		s.addMsg("Already equipped.")
		return
	}

	// Drop current equipped weapon by removing it from backpack (prevents duplication).
	if oldID := s.Player.Equipment.Get(SlotWeapon); oldID != 0 {
		if old, ok := s.Backpack.RemoveByID(oldID); ok {
			s.dropWeaponAdjacent(rng, old)
		}
	}

	// Equip by ID (stable even if inventory order changes).
	s.Player.Equipment.Set(SlotWeapon, chosen.ID)
	s.addMsg("Weapon equipped: " + chosen.DisplayName() + " (+" + itoa(chosen.Strength) + ").")
}

func (s *GameSession) dropWeaponAdjacent(rng *rand.Rand, w Item) {
	dirs := Dirs4
	rng.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
	mw := len(s.Level.Tiles[0])
	mh := len(s.Level.Tiles)
	for _, d := range dirs {
		p := s.PlayerPos.Add(d)
		if !inBounds(p, mw, mh) {
			continue
		}
		if s.Level.Tiles[p.Y][p.X].BlocksMovement() {
			continue
		}
		if _, ok := s.enemyIndexAt(p); ok {
			continue
		}
		ri, ok := s.roomAt(p)
		if ok {
			if s.Level.Rooms[ri].FloorLoot == nil {
				s.Level.Rooms[ri].FloorLoot = map[Point]Item{}
			}
			s.Level.Rooms[ri].FloorLoot[p] = w
			return
		}
		// corridor tile drop: attach to nearest room if any, else ignore.
		for rj := range s.Level.Rooms {
			if s.Level.Rooms[rj].Bounds.Contains(p) {
				s.Level.Rooms[rj].FloorLoot[p] = w
				return
			}
		}
	}
}

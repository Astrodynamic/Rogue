package game

import "math/rand"

func (s *GameSession) advanceLevel(rng *rand.Rand) {
	s.LevelDepth++
	if s.LevelDepth > s.Stats.DeepestLevelReached {
		s.Stats.DeepestLevelReached = s.LevelDepth
	}
	if s.LevelDepth > MaxLevels {
		s.addMsg("You escaped the dungeon!")
		return
	}
	gen := NewGenerator(rng)
	cfg := Config{Width: len(s.Level.Tiles[0]), Height: len(s.Level.Tiles)}
	lvl, startIdx, _ := gen.Generate(s.LevelDepth, cfg)
	startPos := randomFloorInRoom(rng, lvl.Tiles, lvl.Rooms[startIdx].Bounds)

	s.Level = lvl
	s.PlayerPos = startPos
	s.Enemies = nil
	s.Explored = make2D(len(lvl.Tiles[0]), len(lvl.Tiles), false)
	s.populate(rng)
	s.refreshFog()
	s.addMsg("Level " + itoa(s.LevelDepth))
}

func (s *GameSession) populate(rng *rand.Rand) {
	// Scaling rules (simplified): deeper -> more enemies, fewer useful items.
	depth := s.LevelDepth
	enemies := 2 + depth/3
	items := 14 - depth/4
	if items < 2 {
		items = 2
	}

	startRoom := -1
	for i := range s.Level.Rooms {
		if s.Level.Rooms[i].IsStart {
			startRoom = i
			break
		}
	}

	// Place enemies.
	for placed := 0; placed < enemies; {
		ri := rng.Intn(len(s.Level.Rooms))
		if ri == startRoom {
			continue
		}
		p := randomFloorInRoom(rng, s.Level.Tiles, s.Level.Rooms[ri].Bounds)
		if p == s.PlayerPos {
			continue
		}
		if _, ok := s.enemyIndexAt(p); ok {
			continue
		}
		e := makeEnemy(rng, depth)
		e.Pos = p
		s.Enemies = append(s.Enemies, e)
		placed++
	}

	// Place items (no treasures).
	for placed := 0; placed < items; {
		ri := rng.Intn(len(s.Level.Rooms))
		if ri == startRoom {
			continue
		}
		p := randomFloorInRoom(rng, s.Level.Tiles, s.Level.Rooms[ri].Bounds)
		if p == s.PlayerPos || s.Level.Tiles[p.Y][p.X] == TileExit {
			continue
		}
		if _, ok := s.enemyIndexAt(p); ok {
			continue
		}
		if _, ok := s.Level.Rooms[ri].FloorLoot[p]; ok {
			continue
		}
		it := makeItem(rng, depth)
		it.ID = s.nextItemID()
		s.Level.Rooms[ri].FloorLoot[p] = it
		placed++
	}
}

func (s *GameSession) nextItemID() int64 {
	if s.NextItemID < 1 {
		s.NextItemID = 1
	}
	id := s.NextItemID
	s.NextItemID++
	return id
}

func makeEnemy(rng *rand.Rand, depth int) Enemy {
	// Depth scales stats; pick types with rough traits.
	t := EnemyType(rng.Intn(5))
	// Ogre is intentionally dangerous; don't spawn it on the first levels.
	if depth < 4 && t == EnemyOgre {
		t = EnemyZombie
	}
	switch t {
	case EnemyZombie:
		return Enemy{Type: t, Health: 16 + depth, Dexterity: 4 + depth/12, Strength: 7 + depth/3, Hostility: 6 + depth/6}
	case EnemyVampire:
		return Enemy{Type: t, Health: 18 + depth, Dexterity: 10 + depth/5, Strength: 7 + depth/5, Hostility: 8 + depth/6, VampireFirstHitAlwaysMiss: true}
	case EnemyGhost:
		return Enemy{Type: t, Health: 9 + depth/3, Dexterity: 12 + depth/6, Strength: 4 + depth/6, Hostility: 5 + depth/8, GhostInvisible: false}
	case EnemyOgre:
		return Enemy{Type: t, Health: 22 + depth*2, Dexterity: 3 + depth/12, Strength: 10 + depth/2, Hostility: 6 + depth/10}
	case EnemySnakeMage:
		e := Enemy{Type: t, Health: 14 + depth, Dexterity: 14 + depth/2, Strength: 7 + depth/3, Hostility: 10 + depth/2}
		if rng.Intn(2) == 0 {
			e.SnakeDx = 1
		} else {
			e.SnakeDx = -1
		}
		if rng.Intn(2) == 0 {
			e.SnakeDy = 1
		} else {
			e.SnakeDy = -1
		}
		return e
	default:
		return Enemy{Type: EnemyZombie, Health: 18 + depth, Dexterity: 4, Strength: 8 + depth/2, Hostility: 6}
	}
}

func makeItem(rng *rand.Rand, depth int) Item {
	// Deeper: fewer good items; still keep some variety.
	roll := rng.Intn(100)
	if roll < 62 {
		// Food more common.
		return Item{Type: ItemFood, Health: 10 + rng.Intn(12), Description: "Food"}
	}
	if roll < 82 {
		// Elixir (temporary)
		sub := []ItemSubtype{SubtypeDexterity, SubtypeStrength, SubtypeMaxHealth}[rng.Intn(3)]
		it := Item{Type: ItemElixir, Subtype: sub, Temporary: true, Duration: 12 + rng.Intn(10), Description: "Elixir"}
		switch sub {
		case SubtypeDexterity:
			it.Dexterity = 1 + rng.Intn(2)
		case SubtypeStrength:
			it.Strength = 1 + rng.Intn(2)
		case SubtypeMaxHealth:
			it.MaxHealth = 1 + rng.Intn(2)
		}
		return it
	}
	if roll < 90-depth/2 { // fewer scrolls later
		sub := []ItemSubtype{SubtypeDexterity, SubtypeStrength, SubtypeMaxHealth}[rng.Intn(3)]
		it := Item{Type: ItemScroll, Subtype: sub, Description: "Scroll"}
		switch sub {
		case SubtypeDexterity:
			it.Dexterity = 1
		case SubtypeStrength:
			it.Strength = 1
		case SubtypeMaxHealth:
			it.MaxHealth = 1
		}
		return it
	}
	// Weapon
	return Item{Type: ItemWeapon, Strength: 1 + rng.Intn(2) + depth/8, Description: "Sword"}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

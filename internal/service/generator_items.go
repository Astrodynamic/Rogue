package service

import "rogue/internal/domain"

func (g *Generator) GenerateItems(level *domain.Level, depth int, startRoom *domain.Room) {
	g.GenerateItemsWithDifficulty(level, depth, startRoom, 1.0)
}

func (g *Generator) GenerateItemsWithDifficulty(level *domain.Level, depth int, startRoom *domain.Room, difficultyFactor float64) {
	level.Items = make(map[domain.Point]domain.Item)

	itemsPerRoom := domain.ItemsPerRoomBase - depth/7
	if itemsPerRoom < 2 {
		itemsPerRoom = 2
	}
	if difficultyFactor < 1.0 {
		itemsPerRoom += 1
	}

	for _, room := range level.Rooms {
		if startRoom != nil && room.X == startRoom.X && room.Y == startRoom.Y {
			continue
		}

		for i := 0; i < itemsPerRoom; i++ {
			spawnChance := domain.ItemSpawnChance
			if difficultyFactor < 1.0 {
				spawnChance += 10
			}

			if g.rng.IntN(100) < spawnChance {
				item := g.generateRandomItemWithDifficulty(depth, difficultyFactor)
				if item != nil {
					pos := g.getRandomFloorPoint(level, room)
					if pos.X >= 0 && pos.Y >= 0 {
						level.AddItem(pos, item)
					}
				}
			}
		}
	}
}

func (g *Generator) getRandomFloorPoint(level *domain.Level, room domain.Room) domain.Point {
	maxAttempts := 20
	for i := 0; i < maxAttempts; i++ {
		x := g.rng.IntN(room.W-2) + room.X + 1
		y := g.rng.IntN(room.H-2) + room.Y + 1
		p := domain.Point{X: x, Y: y}

		if level.Tiles[y][x].Kind == domain.TileFloor && level.GetItem(p) == nil {
			return p
		}
	}
	return domain.Point{X: -1, Y: -1}
}

func (g *Generator) GenerateTreasureFromEnemyWithDepth(enemy domain.Enemy, depth int) *domain.Treasure {
	if enemy == nil {
		return nil
	}
	actor := enemy.GetActor()
	if actor == nil {
		return nil
	}

	statValue := enemy.GetHostility() + actor.GetStrength() + actor.GetDexterity() + actor.GetMaxHealth()
	value := domain.TreasureBaseValue + (statValue * domain.TreasureStatMultiplier) + (depth * domain.TreasureDepthBonus)

	return &domain.Treasure{
		Value: value,
	}
}

func (g *Generator) generateRandomItem(depth int) domain.Item {
	return g.generateRandomItemWithDifficulty(depth, 1.0)
}

func (g *Generator) generateRandomItemWithDifficulty(depth int, difficultyFactor float64) domain.Item {
	roll := g.rng.IntN(100)

	foodThreshold := domain.FoodSpawnWeight
	elixirThreshold := foodThreshold + 20
	scrollThreshold := elixirThreshold + 15
	treasureThreshold := scrollThreshold + 10
	weaponThreshold := treasureThreshold + 10

	if difficultyFactor < 1.0 {
		foodThreshold += 10
		elixirThreshold += 10
	}

	switch {
	case roll < foodThreshold:
		return &domain.Food{
			Health: domain.FoodBaseHealth + depth*domain.FoodHealthPerDepth,
		}
	case roll < elixirThreshold:
		elixirType := domain.ElixirKind(g.rng.IntN(4))
		return &domain.Elixir{
			ElixirKind: elixirType,
			Amount:     5 + depth,
			Duration:   8 + depth/2,
		}
	case roll < scrollThreshold:
		scrollType := domain.ScrollKind(g.rng.IntN(5))
		amount := 2 + depth/2
		if scrollType == domain.ScrollRegeneration {
			amount = 2 + depth/4
		}
		return &domain.Scroll{
			ScrollKind: scrollType,
			Amount:     amount,
		}
	case roll < treasureThreshold:
		return &domain.Treasure{
			Value: 10 + depth*5,
		}
	case roll < weaponThreshold:
		return &domain.Weapon{
			Strength: 2 + depth/3,
		}
	default:
		armorParts := []domain.ActorPart{domain.ActorPartHead, domain.ActorPartBody, domain.ActorPartLegs}
		partIndex := g.rng.IntN(len(armorParts))
		part := armorParts[partIndex]

		stats := domain.Stats{}
		statRoll := g.rng.IntN(3)
		switch statRoll {
		case 0:
			stats.Strength = 1 + depth/4
		case 1:
			stats.Dexterity = 1 + depth/4
		case 2:
			stats.MaxHealth = 5 + depth/2
		}

		return &domain.Armor{
			ArmorPart: part,
			Stats:     stats,
		}
	}
}

package service

import "rogue/internal/domain"

func (g *Generator) GenerateItems(level *domain.Level, depth int, startRoom *domain.Room) {
	g.GenerateItemsWithDifficulty(level, depth, startRoom, 1.0)
}

func (g *Generator) GenerateItemsWithDifficulty(level *domain.Level, depth int, startRoom *domain.Room, difficultyFactor float64) {
	level.Items = make(map[domain.Point]domain.Item)

	itemsPerRoom := domain.ItemGen.ItemsPerRoom - depth/domain.ItemDepthDivisor
	if itemsPerRoom < domain.ItemMinPerRoom {
		itemsPerRoom = domain.ItemMinPerRoom
	}
	if difficultyFactor < 1.0 {
		itemsPerRoom++
	}

	for _, room := range level.Rooms {
		if startRoom != nil && room.X == startRoom.X && room.Y == startRoom.Y {
			continue
		}

		for i := 0; i < itemsPerRoom; i++ {
			spawnChance := domain.ItemGen.ItemSpawnCh
			if difficultyFactor < 1.0 {
				spawnChance += domain.ItemDifficultyBonus
			}

			if g.rng.IntN(domain.Combat.PercentBase) < spawnChance {
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
	return g.findRandomPositionInRoom(level, room, func(pos domain.Point) bool {
		return level.Tiles[pos.Y][pos.X].Kind == domain.TileFloor && level.GetItem(pos) == nil
	})
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
	value := domain.TreasureBaseVal + (statValue * domain.TreasureStatMult) + (depth * domain.TreasureDepthBon)

	return &domain.Treasure{
		Value: value,
	}
}

func (g *Generator) generateRandomItem(depth int) domain.Item {
	return g.generateRandomItemWithDifficulty(depth, 1.0)
}

func (g *Generator) generateRandomItemWithDifficulty(depth int, difficultyFactor float64) domain.Item {
	roll := g.rng.IntN(domain.Combat.PercentBase)

	foodThreshold := domain.ItemGen.FoodSpawnWgt
	elixirThreshold := foodThreshold + domain.ItemElixirThresholdOffset
	scrollThreshold := elixirThreshold + domain.ItemScrollThresholdOffset
	treasureThreshold := scrollThreshold + domain.ItemTreasureThresholdOffset
	weaponThreshold := treasureThreshold + domain.ItemWeaponThresholdOffset

	if difficultyFactor < 1.0 {
		foodThreshold += domain.ItemDifficultyBonus
		elixirThreshold += domain.ItemDifficultyBonus
	}

	switch {
	case roll < foodThreshold:
		return &domain.Food{
			Health: domain.ItemGen.FoodBaseHP + depth*domain.ItemGen.FoodHPPerDepth,
		}
	case roll < elixirThreshold:
		elixirType := domain.ElixirKind(g.rng.IntN(domain.ElixirKindCount))
		return &domain.Elixir{
			ElixirKind: elixirType,
			Amount:     domain.ElixirBaseAmount + depth,
			Duration:   domain.ElixirBaseDuration + depth/domain.ElixirDurationDivisor,
		}
	case roll < scrollThreshold:
		scrollType := domain.ScrollKind(g.rng.IntN(domain.ScrollKindCount))
		amount := domain.ScrollBaseAmount + depth/domain.ScrollAmountDivisor
		if scrollType == domain.ScrollRegeneration {
			amount = domain.ScrollBaseAmount + depth/domain.ScrollRegenDivisor
		}
		return &domain.Scroll{
			ScrollKind: scrollType,
			Amount:     amount,
		}
	case roll < treasureThreshold:
		return &domain.Treasure{
			Value: domain.TreasureBaseDropValue + depth*domain.TreasureDepthMultiplier,
		}
	case roll < weaponThreshold:
		return &domain.Weapon{
			Strength: domain.WeaponBaseStrength + depth/domain.WeaponStrengthDivisor,
		}
	default:
		armorParts := []domain.ActorPart{domain.ActorPartHead, domain.ActorPartBody, domain.ActorPartLegs}
		partIndex := g.rng.IntN(len(armorParts))
		part := armorParts[partIndex]

		stats := domain.Stats{}
		statRoll := g.rng.IntN(domain.ArmorStatRollCount)
		switch statRoll {
		case 0:
			stats.Strength = domain.ArmorStatBase + depth/domain.ArmorStatDivisor
		case 1:
			stats.Dexterity = domain.ArmorStatBase + depth/domain.ArmorStatDivisor
		case 2:
			stats.MaxHealth = domain.ArmorHealthBase + depth/domain.ArmorHealthDivisor
		}

		return &domain.Armor{
			ArmorPart: part,
			Stats:     stats,
		}
	}
}

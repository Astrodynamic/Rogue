package service

import "rogue/internal/domain"

func (g *Generator) GenerateItems(level *domain.Level, depth int, startRoom *domain.Room) {
	level.Items = make(map[domain.Point]domain.Item)

	itemsPerRoom := 1 + depth/3
	if itemsPerRoom > 3 {
		itemsPerRoom = 3
	}

	for _, room := range level.Rooms {
		if startRoom != nil && room.X == startRoom.X && room.Y == startRoom.Y {
			continue
		}

		for i := 0; i < itemsPerRoom; i++ {
			if g.rng.IntN(100) < 60 {
				item := g.generateRandomItem(depth)
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
	roll := g.rng.IntN(100)

	switch {
	case roll < 20:
		return &domain.Food{
			Health: 10 + depth*2,
		}
	case roll < 35:
		elixirType := domain.ElixirKind(g.rng.IntN(4))
		return &domain.Elixir{
			ElixirKind: elixirType,
			Amount:     5 + depth,
			Duration:   5 + depth/2,
		}
	case roll < 50:
		scrollType := domain.ScrollKind(g.rng.IntN(4))
		return &domain.Scroll{
			ScrollKind: scrollType,
			Amount:     3 + depth,
		}
	case roll < 60:
		return &domain.Treasure{
			Value: 10 + depth*5,
		}
	case roll < 75:
		return &domain.Weapon{
			Strength: 2 + depth/2,
		}
	case roll < 90:
		armorParts := []domain.ActorPart{domain.ActorPartHead, domain.ActorPartBody, domain.ActorPartLegs}
		partIndex := g.rng.IntN(len(armorParts))
		part := armorParts[partIndex]

		stats := domain.Stats{}
		statRoll := g.rng.IntN(3)
		switch statRoll {
		case 0:
			stats.Strength = 1 + depth/3
		case 1:
			stats.Dexterity = 1 + depth/3
		case 2:
			stats.MaxHealth = 2 + depth/2
		}

		return &domain.Armor{
			ArmorPart: part,
			Stats:     stats,
		}
	default:
		return nil
	}
}

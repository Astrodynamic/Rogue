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

func (g *Generator) generateRandomItem(depth int) domain.Item {
	roll := g.rng.IntN(100)

	switch {
	case roll < 20:
		return &domain.Food{
			Health: 10 + depth*2,
		}
	case roll < 40:
		elixirType := domain.ElixirKind(g.rng.IntN(4))
		return &domain.Elixir{
			ElixirKind: elixirType,
			Amount:     5 + depth,
			Duration:   5 + depth/2,
		}
	case roll < 60:
		scrollType := domain.ScrollKind(g.rng.IntN(4))
		return &domain.Scroll{
			ScrollKind: scrollType,
			Amount:     3 + depth,
		}
	case roll < 75:
		return &domain.Treasure{
			Value: 10 + depth*5,
		}
	case roll < 90:
		return &domain.Weapon{
			Strength: 2 + depth/2,
		}
	default:
		return nil
	}
}

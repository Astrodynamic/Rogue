package service

import (
	"rogue/internal/domain"
)

func (g *Generator) GenerateEnemies(level *domain.Level, depth int, startRoom *domain.Room) {
	g.GenerateEnemiesWithDifficulty(level, depth, startRoom, 1.0)
}

func (g *Generator) GenerateEnemiesWithDifficulty(level *domain.Level, depth int, startRoom *domain.Room, difficultyFactor float64) {
	level.Enemies = make(map[domain.Point]domain.Enemy)

	enemiesPerRoomFloat := (float64(domain.EnemyGeneration.CountBase) + float64(depth)*domain.EnemyGeneration.CountScaling) * difficultyFactor
	enemiesPerRoom := int(enemiesPerRoomFloat)
	if enemiesPerRoom > domain.EnemyGeneration.MaxCountPerRoom {
		enemiesPerRoom = domain.EnemyGeneration.MaxCountPerRoom
	}
	if enemiesPerRoom < 1 {
		enemiesPerRoom = 1
	}

	for _, room := range level.Rooms {
		if startRoom != nil && room.X == startRoom.X && room.Y == startRoom.Y {
			continue
		}

		for i := 0; i < enemiesPerRoom; i++ {
			spawnChance := domain.EnemyGeneration.SpawnChance
			if difficultyFactor > 1.0 {
				spawnChance += domain.EnemySpawnChanceDifficultyBonus
			} else if difficultyFactor < 1.0 {
				spawnChance -= domain.EnemySpawnChanceDifficultyPenalty
			}

			if g.rng.IntN(domain.Combat.PercentBase) < spawnChance {
				enemy := g.generateRandomEnemy(depth)
				if enemy != nil {
					pos := g.getRandomEnemyPoint(level, room)
					if pos.X >= 0 && pos.Y >= 0 {
						level.AddEnemy(pos, enemy)
					}
				}
			}
		}
	}
}

func (g *Generator) generateRandomEnemy(depth int) domain.Enemy {
	var possibleEnemies []func(int) domain.Enemy

	if depth <= domain.ZombieMaxDepth {
		possibleEnemies = append(possibleEnemies, func(d int) domain.Enemy {
			return domain.NewZombie(d)
		})
	}

	if depth >= domain.GhostMinDepth {
		possibleEnemies = append(possibleEnemies, func(d int) domain.Enemy {
			return domain.NewGhost(d)
		})
	}

	if depth >= domain.VampireMinDepth {
		possibleEnemies = append(possibleEnemies, func(d int) domain.Enemy {
			return domain.NewVampire(d)
		})
	}

	if depth >= domain.OgreMinDepth {
		possibleEnemies = append(possibleEnemies, func(d int) domain.Enemy {
			return domain.NewOgre(d)
		})
	}

	if depth >= domain.SnakeMageMinDepth {
		possibleEnemies = append(possibleEnemies, func(d int) domain.Enemy {
			return domain.NewSnakeMage(d)
		})
	}

	if depth >= domain.MimicMinDepth {
		possibleEnemies = append(possibleEnemies, func(d int) domain.Enemy {

			disguises := []domain.ItemKind{
				domain.ItemFood,
				domain.ItemElixir,
				domain.ItemScroll,
				domain.ItemWeapon,
				domain.ItemTreasure,
			}
			disguise := disguises[g.rng.IntN(len(disguises))]
			return domain.NewMimic(d, disguise)
		})
	}

	if len(possibleEnemies) == 0 {
		return domain.NewZombie(depth)
	}

	weights := make([]int, len(possibleEnemies))
	for i := range possibleEnemies {
		weights[i] = depth + 1
	}

	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}

	roll := g.rng.IntN(totalWeight)
	currentWeight := 0
	for i, weight := range weights {
		currentWeight += weight
		if roll < currentWeight {
			return possibleEnemies[i](depth)
		}
	}

	return possibleEnemies[0](depth)
}

func (g *Generator) getRandomEnemyPoint(level *domain.Level, room domain.Room) domain.Point {
	return g.findRandomPositionInRoom(level, room, func(pos domain.Point) bool {
		return level.Tiles[pos.Y][pos.X].Kind == domain.TileFloor &&
			level.GetItem(pos) == nil && level.GetEnemy(pos) == nil
	})
}

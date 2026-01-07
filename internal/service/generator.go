package service

import (
	"rogue/internal/domain"
)

type Generator struct {
	rng domain.RandomGenerator
}

func NewGeneratorWithRNG(rng domain.RandomGenerator) *Generator {
	return &Generator{
		rng: rng,
	}
}

func (g *Generator) Generate(world *domain.World) {
	g.GenerateLevel(world.Level)
	room := g.GenerateStartPosition(world)
	depth := world.GetDepth()
	difficultyFactor := world.GameState.DifficultyFactor

	g.GenerateItemsWithDifficulty(world.Level, depth, room, difficultyFactor)
	g.GenerateEnemiesWithDifficulty(world.Level, depth, room, difficultyFactor)
}

func (g *Generator) GenerateStartPosition(world *domain.World) *domain.Room {
	room := world.Level.Rooms[g.rng.IntN(len(world.Level.Rooms))]

	maxAttempts := 50
	for i := 0; i < maxAttempts; i++ {
		x := g.rng.IntN(room.W-2) + room.X + 1
		y := g.rng.IntN(room.H-2) + room.Y + 1
		p := domain.Point{X: x, Y: y}

		if world.Level.Tiles[y][x].Kind == domain.TileFloor {
			world.Player.Point = p
			return &room
		}
	}

	world.Player.Point = room.Center()
	return &room
}

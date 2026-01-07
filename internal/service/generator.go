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

func (g *Generator) findRandomPositionInRoom(level *domain.Level, room domain.Room, validator func(domain.Point) bool) domain.Point {
	maxAttempts := domain.LevelGeneration.MaxRandomPositionAttempts
	for i := 0; i < maxAttempts; i++ {
		x := g.rng.IntN(room.W-domain.LevelGeneration.RoomPadding) + room.X + 1
		y := g.rng.IntN(room.H-domain.LevelGeneration.RoomPadding) + room.Y + 1
		p := domain.Point{X: x, Y: y}

		if validator(p) {
			return p
		}
	}
	return domain.Point{X: -1, Y: -1}
}

func (g *Generator) GenerateStartPosition(world *domain.World) *domain.Room {
	room := world.Level.Rooms[g.rng.IntN(len(world.Level.Rooms))]

	p := g.findRandomPositionInRoom(world.Level, room, func(pos domain.Point) bool {
		return world.Level.Tiles[pos.Y][pos.X].Kind == domain.TileFloor
	})

	if p.X >= 0 {
		world.Player.Point = p
	} else {
		world.Player.Point = room.Center()
	}
	return &room
}

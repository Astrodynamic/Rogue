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
	g.GenerateItems(world.Level, depth, room)
	g.GenerateEnemies(world.Level, depth, room)
	world.GameState.AdvanceLevel()
}

func (g *Generator) GenerateStartPosition(world *domain.World) *domain.Room {
	room := world.Level.Rooms[g.rng.IntN(len(world.Level.Rooms))]
	world.Player.Point = room.Center()
	return &room
}

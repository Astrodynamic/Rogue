package service

import (
	"math/rand/v2"
	"time"

	"rogue/internal/domain"
)

type Generator struct {
	rng *rand.Rand
}

func NewGenerator() *Generator {
	seed := uint64(time.Now().UnixNano())
	return &Generator{
		rng: rand.New(rand.NewPCG(seed, seed)),
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

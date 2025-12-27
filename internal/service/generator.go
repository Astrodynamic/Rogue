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
	g.GenerateStartPosition(world)
}

func (g *Generator) GenerateStartPosition(world *domain.World) {
	room := world.Level.Rooms[g.rng.IntN(len(world.Level.Rooms))]
	world.Player.Point = room.Center()
}

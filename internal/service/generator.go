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
}

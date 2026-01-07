package service

import (
	"math/rand/v2"
	"rogue/internal/domain"
)

type RandAdapter struct {
	rng *rand.Rand
}

func NewRandAdapter(rng *rand.Rand) *RandAdapter {
	return &RandAdapter{rng: rng}
}

func (r *RandAdapter) Float64() float64 {
	return r.rng.Float64()
}

func (r *RandAdapter) IntN(n int) int {
	return r.rng.IntN(n)
}

func (r *RandAdapter) Shuffle(n int, swap func(i, j int)) {
	r.rng.Shuffle(n, swap)
}

var _ domain.RandomGenerator = (*RandAdapter)(nil)

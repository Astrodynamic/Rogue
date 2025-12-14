package game

import "math/rand"

func HitChance(attDex, defDex int) float64 {
	// Simple bounded formula; tune later.
	// p = (attDex + 10) / (attDex + defDex + 20)
	num := float64(attDex + 10)
	den := float64(attDex + defDex + 20)
	if den <= 0 {
		return 0.5
	}
	p := num / den
	if p < 0.05 {
		return 0.05
	}
	if p > 0.95 {
		return 0.95
	}
	return p
}

func RollHit(rng *rand.Rand, attDex, defDex int) bool {
	return rng.Float64() < HitChance(attDex, defDex)
}

func Damage(rng *rand.Rand, strength int, weaponBonus int) int {
	base := 1 + strength/2
	if base < 1 {
		base = 1
	}
	mod := weaponBonus
	// Small randomness.
	d := base + mod + rng.Intn(3) // +0..2
	if d < 1 {
		d = 1
	}
	return d
}

func TreasureDrop(rng *rand.Rand, e Enemy) int {
	// Difficulty-derived drop; grows with depth through enemy stats.
	score := e.Hostility + e.Strength*2 + e.Dexterity + e.Health/2
	if score < 1 {
		score = 1
	}
	return rng.Intn(score) + score/2
}

func XPDrop(e Enemy) int {
	// XP scales with enemy difficulty; tuned to keep pace by mid depth.
	score := e.Hostility + e.Strength*2 + e.Dexterity + e.Health/2
	if score < 1 {
		score = 1
	}
	// Typical early: ~10-20 XP; mid: ~25-45; late: higher.
	return 6 + score/4
}

package domain

func XPToNextLevel(level int) int {
	if level < 1 {
		level = 1
	}
	// Ramp tuned so the player actually keeps up with depth scaling.
	// L1->2: 30, then grows by +10 per level: 30,40,50,60,...
	return 30 + (level-1)*10
}

// GainXP applies XP, levels up, and returns how many levels were gained.
// Level-ups increase survivability and damage slightly so the player keeps up with depth scaling.
func (c *Character) GainXP(amount int) (levelsGained int) {
	if amount <= 0 {
		return 0
	}
	c.XP += amount
	for c.XP >= XPToNextLevel(c.Level) {
		c.XP -= XPToNextLevel(c.Level)
		c.Level++
		levelsGained++

		// Growth: every level +3 max HP, +1 strength; every 2nd level +1 dex.
		c.MaxHealth += 3
		c.Health += 3
		c.Strength += 1
		if c.Level%2 == 0 {
			c.Dexterity += 1
		}
	}
	return levelsGained
}

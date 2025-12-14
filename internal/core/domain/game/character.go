package game

type Character struct {
	MaxHealth int
	Health    int
	Dexterity int
	Strength  int
	XP        int
	Level     int
	Equipment Equipment
}

func NewCharacter(maxHP, dex, str int) Character {
	if maxHP < 1 {
		maxHP = 1
	}
	return Character{
		MaxHealth: maxHP,
		Health:    maxHP,
		Dexterity: dex,
		Strength:  str,
		XP:        0,
		Level:     1,
		Equipment: Equipment{},
	}
}

func (c *Character) IsDead() bool { return c.Health <= 0 }

func (c *Character) ApplyDamage(dmg int) {
	if dmg <= 0 {
		return
	}
	c.Health -= dmg
}

func (c *Character) Heal(hp int) {
	if hp <= 0 {
		return
	}
	c.Health += hp
	if c.Health > c.MaxHealth {
		c.Health = c.MaxHealth
	}
}

func (c *Character) IncreaseMaxHealth(delta int) {
	if delta <= 0 {
		return
	}
	c.MaxHealth += delta
	c.Health += delta
}

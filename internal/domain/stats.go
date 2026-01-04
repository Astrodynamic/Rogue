package domain

type Stats struct {
	MaxHealth int
	Health    int
	Dexterity int
	Strength  int
}

func (s *Stats) GetMaxHealth() int {
	return s.MaxHealth
}

func (s *Stats) AddMaxHealth(amount int) {
	s.MaxHealth += amount
	if s.MaxHealth < 1 {
		s.MaxHealth = 1
	}

	if amount > 0 {
		s.Health += amount
	}

	if s.Health > s.MaxHealth {
		s.Health = s.MaxHealth
	}
}

func (s *Stats) GetHealth() int {
	return s.Health
}

func (s *Stats) AddHealth(amount int) {
	s.Health += amount
	if s.Health > s.MaxHealth {
		s.Health = s.MaxHealth
	}

	if s.Health < 0 {
		s.Health = 0
	}
}

func (s *Stats) GetDexterity() int {
	return s.Dexterity
}

func (s *Stats) AddDexterity(amount int) {
	s.Dexterity += amount
	if s.Dexterity < 0 {
		s.Dexterity = 0
	}
}

func (s *Stats) GetStrength() int {
	return s.Strength
}

func (s *Stats) AddStrength(amount int) {
	s.Strength += amount
	if s.Strength < 0 {
		s.Strength = 0
	}
}

func (s *Stats) Apply(changes Stats) {
	s.AddMaxHealth(changes.MaxHealth)
	s.AddHealth(changes.Health)
	s.AddDexterity(changes.Dexterity)
	s.AddStrength(changes.Strength)
}

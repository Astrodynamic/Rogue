package domain

type ScrollKind int

const (
	ScrollHealth ScrollKind = iota
	ScrollMaxHealth
	ScrollDexterity
	ScrollStrength
)

type Scroll struct {
	ScrollKind ScrollKind
	Amount     int
}

func (s *Scroll) Type() ItemKind {
	return ItemScroll
}

func (s *Scroll) Name() string {
	return "Scroll"
}

func (s *Scroll) Stackable() bool {
	return true
}

func (s *Scroll) MaxStack() int {
	return 9
}

func (s *Scroll) Equals(other Item) bool {
	if other.Type() != ItemScroll {
		return false
	}
	otherScroll, ok := other.(*Scroll)
	if !ok {
		return false
	}
	return s.ScrollKind == otherScroll.ScrollKind && s.Amount == otherScroll.Amount
}

func (s *Scroll) Use() ItemUseResult {
	var effect Effect
	switch s.ScrollKind {
	case ScrollHealth:
		effect = &HealthEffect{
			BaseEffect: BaseEffect{
				duration: 0,
			},
			Stats: Stats{Health: s.Amount},
		}
	case ScrollMaxHealth:
		effect = &MaxHealthEffect{
			BaseEffect: BaseEffect{
				duration: 0,
			},
			Stats: Stats{MaxHealth: s.Amount},
		}
	case ScrollDexterity:
		effect = &DexterityEffect{
			BaseEffect: BaseEffect{
				duration: 0,
			},
			Stats: Stats{Dexterity: s.Amount},
		}
	case ScrollStrength:
		effect = &StrengthEffect{
			BaseEffect: BaseEffect{
				duration: 0,
			},
			Stats: Stats{Strength: s.Amount},
		}
	}

	return ItemUseResult{
		Success:  true,
		Consumed: true,
		Message:  s.Name() + " used",
		Effects:  []Effect{effect},
	}
}

func (s *Scroll) GetStats() Stats {
	switch s.ScrollKind {
	case ScrollHealth:
		return Stats{Health: s.Amount}
	case ScrollMaxHealth:
		return Stats{MaxHealth: s.Amount}
	case ScrollDexterity:
		return Stats{Dexterity: s.Amount}
	case ScrollStrength:
		return Stats{Strength: s.Amount}
	default:
		return Stats{}
	}
}

package domain

type ScrollKind int

const (
	ScrollHealth ScrollKind = iota
	ScrollMaxHealth
	ScrollDexterity
	ScrollStrength
	ScrollRegeneration
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

func (s *Scroll) createEffect() Effect {
	switch s.ScrollKind {
	case ScrollHealth:
		return NewStatEffect(Stats{Health: s.Amount}, 0)
	case ScrollMaxHealth:
		return NewStatEffect(Stats{MaxHealth: s.Amount}, 0)
	case ScrollDexterity:
		return NewStatEffect(Stats{Dexterity: s.Amount}, 0)
	case ScrollStrength:
		return NewStatEffect(Stats{Strength: s.Amount}, 0)
	case ScrollRegeneration:
		return &RegenerationEffect{
			BaseEffect: BaseEffect{duration: -1},
			Amount:     s.Amount,
		}
	default:
		return nil
	}
}

func (s *Scroll) getEffectStats() Stats {
	switch s.ScrollKind {
	case ScrollHealth, ScrollRegeneration:
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

func (s *Scroll) Use() ItemUseResult {
	effect := s.createEffect()
	if effect == nil {
		return ItemUseResult{
			Success: false,
			Message: "Invalid scroll",
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
	return s.getEffectStats()
}

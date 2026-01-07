package domain

type ElixirKind int

const (
	ElixirHealth ElixirKind = iota
	ElixirMaxHealth
	ElixirDexterity
	ElixirStrength
)

type Elixir struct {
	ElixirKind ElixirKind
	Amount     int
	Duration   int
}

func (e *Elixir) Type() ItemKind {
	return ItemElixir
}

func (e *Elixir) Name() string {
	return "Elixir"
}

func (e *Elixir) Stackable() bool {
	return true
}

func (e *Elixir) MaxStack() int {
	return 9
}

func (e *Elixir) Equals(other Item) bool {
	if other.Type() != ItemElixir {
		return false
	}
	otherElixir, ok := other.(*Elixir)
	if !ok {
		return false
	}
	return e.ElixirKind == otherElixir.ElixirKind && e.Amount == otherElixir.Amount && e.Duration == otherElixir.Duration
}

func (e *Elixir) getEffectStats() Stats {
	switch e.ElixirKind {
	case ElixirHealth:
		return Stats{Health: e.Amount}
	case ElixirMaxHealth:
		return Stats{MaxHealth: e.Amount}
	case ElixirDexterity:
		return Stats{Dexterity: e.Amount}
	case ElixirStrength:
		return Stats{Strength: e.Amount}
	default:
		return Stats{}
	}
}

func (e *Elixir) Use() ItemUseResult {
	stats := e.getEffectStats()
	duration := e.Duration
	if e.ElixirKind == ElixirHealth {
		duration = 0
	}

	return ItemUseResult{
		Success:  true,
		Consumed: true,
		Message:  e.Name() + " consumed",
		Effects:  []Effect{NewStatEffect(stats, duration)},
	}
}

func (e *Elixir) GetStats() Stats {
	return e.getEffectStats()
}

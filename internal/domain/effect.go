package domain

type Effect interface {
	Apply(target *Actor) error
	Revert(target *Actor) error
	Duration() int
	Tick(target *Actor) bool
}

type BaseEffect struct {
	duration int
}

func (e *BaseEffect) Duration() int {
	return e.duration
}

func (e *BaseEffect) Tick(target *Actor) bool {
	if e.duration > 0 {
		e.duration--
		return e.duration == 0
	}
	return false
}

type StatEffect struct {
	BaseEffect
	Stats Stats
}

func NewStatEffect(stats Stats, duration int) *StatEffect {
	return &StatEffect{
		BaseEffect: BaseEffect{duration: duration},
		Stats:      stats,
	}
}

func (e *StatEffect) Apply(target *Actor) error {
	target.Stats.Apply(e.Stats)
	return nil
}

func (e *StatEffect) Revert(target *Actor) error {
	revert := Stats{
		MaxHealth: -e.Stats.MaxHealth,
		Health:    -e.Stats.Health,
		Dexterity: -e.Stats.Dexterity,
		Strength:  -e.Stats.Strength,
	}
	target.Stats.Apply(revert)
	return nil
}

func (e *StatEffect) Tick(target *Actor) bool {
	return e.BaseEffect.Tick(target)
}

func NewMaxHealthEffect(amount int, duration int) *StatEffect {
	return NewStatEffect(Stats{MaxHealth: amount}, duration)
}

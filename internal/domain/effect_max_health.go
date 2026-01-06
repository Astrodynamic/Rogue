package domain

type MaxHealthEffect struct {
	BaseEffect
	Stats Stats
}

func (e *MaxHealthEffect) Apply(target *Actor) error {
	target.Stats.Apply(e.Stats)
	return nil
}

func (e *MaxHealthEffect) Revert(target *Actor) error {
	revert := Stats{
		MaxHealth: -e.Stats.MaxHealth,
	}
	target.Stats.Apply(revert)
	return nil
}

func (e *MaxHealthEffect) Tick(target *Actor) bool {
	return e.BaseEffect.Tick(target)
}

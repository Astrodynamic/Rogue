package domain

type HealthEffect struct {
	BaseEffect
	Stats Stats
}

func (e *HealthEffect) Apply(target *Actor) error {
	target.Stats.Apply(e.Stats)
	return nil
}

func (e *HealthEffect) Revert(target *Actor) error {
	revert := Stats{
		Health: -e.Stats.Health,
	}
	target.Stats.Apply(revert)
	return nil
}

func (e *HealthEffect) Tick(target *Actor) bool {
	return e.BaseEffect.Tick(target)
}

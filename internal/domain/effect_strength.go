package domain

type StrengthEffect struct {
	BaseEffect
	Stats Stats
}

func (e *StrengthEffect) Apply(target *Actor) error {
	target.Stats.Apply(e.Stats)
	return nil
}

func (e *StrengthEffect) Revert(target *Actor) error {
	revert := Stats{
		Strength: -e.Stats.Strength,
	}
	target.Stats.Apply(revert)
	return nil
}

package domain

type DexterityEffect struct {
	BaseEffect
	Stats Stats
}

func (e *DexterityEffect) Apply(target *Actor) error {
	target.Stats.Apply(e.Stats)
	return nil
}

func (e *DexterityEffect) Revert(target *Actor) error {
	revert := Stats{
		Dexterity: -e.Stats.Dexterity,
	}
	target.Stats.Apply(revert)
	return nil
}

func (e *DexterityEffect) Tick(target *Actor) bool {
	return e.BaseEffect.Tick(target)
}

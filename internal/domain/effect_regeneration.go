package domain

type RegenEffect struct {
	BaseEffect
	Amount int
}

func (e *RegenEffect) Apply(target *Actor) error {
	return nil
}

func (e *RegenEffect) Revert(target *Actor) error {
	return nil
}

func (e *RegenEffect) Tick(target *Actor) bool {

	target.AddHealth(e.Amount)

	if e.duration < 0 {
		return false
	}
	return e.BaseEffect.Tick(target)
}

package domain

type RegenerationEffect struct {
	BaseEffect
	Amount int
}

func (e *RegenerationEffect) Apply(target *Actor) error {
	return nil
}

func (e *RegenerationEffect) Revert(target *Actor) error {
	return nil
}

func (e *RegenerationEffect) Tick(target *Actor) bool {

	target.AddHealth(e.Amount)

	if e.duration < 0 {
		return false
	}
	return e.BaseEffect.Tick(target)
}

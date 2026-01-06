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
	if e.duration > 0 {
		target.AddHealth(e.Amount)
	}
	return e.BaseEffect.Tick(target)
}

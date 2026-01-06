package domain

type ActorPart int

const (
	ActorPartHead ActorPart = iota
	ActorPartBody
	ActorPartHand
	ActorPartLegs
)

type ActorState int

const (
	ActorStateNormal ActorState = iota
	ActorStateSleep
)

type Actor struct {
	Point
	Stats
	Name     string
	State    ActorState
	Effects  []Effect
	Backpack *Backpack
}

func (a *Actor) Move(dir Point) {
	if a.State == ActorStateSleep {
		return
	}
	a.Point = a.Point.Add(dir)
}

func (a *Actor) Attack(target *Actor) {

}

func (a *Actor) AddEffect(effect Effect) error {
	if err := effect.Apply(a); err != nil {
		return err
	}

	if effect.Duration() > 0 {
		a.Effects = append(a.Effects, effect)
	}

	return nil
}

func (a *Actor) RemoveEffect(effect Effect) error {
	if err := effect.Revert(a); err != nil {
		return err
	}

	for i, e := range a.Effects {
		if e == effect {
			a.Effects = append(a.Effects[:i], a.Effects[i+1:]...)
			break
		}
	}

	return nil
}

func (a *Actor) TickEffects() {
	var expired []Effect

	for _, effect := range a.Effects {
		if effect.Tick(a) {
			expired = append(expired, effect)
		}
	}

	for _, effect := range expired {
		a.RemoveEffect(effect)
	}
}

func (a *Actor) UseItem(item Item) ItemUseResult {
	result := item.Use()
	if !result.Success {
		return result
	}

	for _, effect := range result.Effects {
		if err := a.AddEffect(effect); err != nil {
			return ItemUseResult{
				Success:  false,
				Consumed: false,
				Error:    err,
				Message:  "Failed to apply effect",
			}
		}
	}

	return result
}

func (a *Actor) GetEffectiveStats() Stats {
	return a.Stats
}

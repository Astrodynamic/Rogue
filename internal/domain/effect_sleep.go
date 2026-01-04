package domain

type SleepEffect struct {
	BaseEffect
}

func (e *SleepEffect) Apply(target *Actor) error {
	target.State = ActorStateSleep
	return nil
}

func (e *SleepEffect) Revert(target *Actor) error {
	target.State = ActorStateNormal
	return nil
}

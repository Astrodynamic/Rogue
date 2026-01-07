package domain

type SleepEffect struct {
	BaseEffect
}

func NewSleepEffect(duration int) *SleepEffect {
	return &SleepEffect{
		BaseEffect: BaseEffect{duration: duration},
	}
}

func (e *SleepEffect) Apply(target *Actor) error {
	target.State = ActorStateSleep
	return nil
}

func (e *SleepEffect) Revert(target *Actor) error {
	target.State = ActorStateNormal
	return nil
}

func (e *SleepEffect) Tick(target *Actor) bool {
	return e.BaseEffect.Tick(target)
}

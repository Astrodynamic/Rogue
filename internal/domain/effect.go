package domain

type Effect interface {
	Apply(target *Actor) error
	Revert(target *Actor) error
	Duration() int
	Tick() bool
}

type BaseEffect struct {
	duration int
}

func (e *BaseEffect) Duration() int {
	return e.duration
}

func (e *BaseEffect) Tick() bool {
	if e.duration > 0 {
		e.duration--
		return e.duration == 0
	}
	return false
}

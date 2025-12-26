package domain

type ActorPart int

const (
	ActorPartHead ActorPart = iota
	ActorPartBody
	ActorPartHand
	ActorPartLegs
)

type Actor struct {
	Point
	Stats
	Effects []Effect
	Name    string
}

func (a *Actor) Move(dir Point) {
	a.Point = a.Point.Add(dir)
}

func (a *Actor) Attack(target *Actor) {

}

package domain

type Player struct {
	Actor
	Equipment *Equipment
}

func NewPlayer() *Player {
	return &Player{
		Actor: Actor{
			Point:    Point{X: 0, Y: 0},
			Backpack: NewBackpack(),
		},
		Equipment: NewEquipment(),
	}
}

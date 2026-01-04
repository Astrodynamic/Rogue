package domain

type Enemy struct {
	Actor
	Hostility int
}

func NewEnemy() *Enemy {
	return &Enemy{
		Actor: Actor{
			Backpack: NewBackpack(),
		},
		Hostility: 0,
	}
}

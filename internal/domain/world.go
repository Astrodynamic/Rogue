package domain

type World struct {
	Level  *Level
	Player *Player
	Depth  int
}

func NewWorld(width, height int) *World {
	return &World{
		Level:  NewLevel(width, height),
		Player: NewPlayer(),
		Depth:  1,
	}
}

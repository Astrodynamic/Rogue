package domain

type World struct {
	Level  *Level
	Player *Player
}

func NewWorld(width, height int) *World {
	return &World{
		Level:  NewLevel(width, height),
		Player: NewPlayer(),
	}
}

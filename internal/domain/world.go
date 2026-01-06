package domain

type World struct {
	Level     *Level
	Player    *Player
	GameState *GameState
}

func NewWorld(width, height int) *World {
	return &World{
		Level:     NewLevel(width, height),
		Player:    NewPlayer(),
		GameState: NewGameState(),
	}
}

func (w *World) GetDepth() int {
	return w.GameState.Depth
}

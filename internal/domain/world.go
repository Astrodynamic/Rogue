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

func NewWorldWithPlayerName(width, height int, playerName string) *World {
	world := NewWorld(width, height)
	world.GameState.PlayerName = playerName
	world.Player.Name = playerName
	return world
}

func (w *World) GetDepth() int {
	return w.GameState.Depth
}

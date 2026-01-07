package service

import (
	"rogue/internal/domain"
)

func (g *Game) loadGame() {
	world, err := g.LoadGameState()
	if err == nil && world != nil {
		g.World = world
		g.UpdateVisibility()
		g.ui.AddLog("Game loaded successfully!")
	} else {
		g.ui.AddLog("Failed to load game. Starting new game...")
		g.startNewGame()
	}
}

func (g *Game) startNewGame() {
	world := domain.NewWorld(domain.Width, domain.Height)
	g.generator.Generate(world)
	g.World = world
	g.UpdateVisibility()
}

func (g *Game) handleGameCompletion() {
	g.SaveStatistics()
	g.ui.AddLog("Congratulations! You have completed all 21 levels!")
	g.ui.AddLog("Game completed. Starting new game...")
	world := domain.NewWorld(domain.Width, domain.Height)
	g.World = world
	g.generator.Generate(g.World)
	g.UpdateVisibility()
}

func (g *Game) SaveStatistics() error {
	playthrough := g.World.GameState.ToPlaythroughStatistics()
	return g.statisticsStore.SavePlaythrough(playthrough)
}

func (g *Game) SaveGameState() error {
	return g.gameStateStore.Save(g.World)
}

func (g *Game) LoadGameState() (*domain.World, error) {
	return g.gameStateStore.Load()
}

func (g *Game) HasSaveGame() bool {
	return g.gameStateStore.HasSave()
}

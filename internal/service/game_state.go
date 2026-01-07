package service

import (
	"fmt"

	"rogue/internal/domain"
)

func (g *Game) loadGame() {
	saves, err := g.gameStateStore.ListSaves()
	if err != nil || len(saves) == 0 {
		g.ui.AddLog("No saved games found. Starting new game...")
		g.startNewGame()
		return
	}

	playerName := saves[0]
	if len(saves) > 1 {
		playerName = saves[0]
	}

	world, err := g.LoadGameState(playerName)
	if err == nil && world != nil {
		g.World = world
		g.UpdateVisibility()
		g.ui.AddLog("Game loaded successfully for " + playerName + "!")
	} else {
		g.ui.AddLog("Failed to load game. Starting new game...")
		g.startNewGame()
	}
}

func (g *Game) startNewGame() {
	g.ui.SetNameInputState(true)
}

func (g *Game) handleGameCompletion() {
	g.SaveStatistics()
	g.ui.AddLog("Congratulations! You have completed all 21 levels!")
	g.ui.AddLog("Game completed! Final score: " + fmt.Sprintf("%d", g.World.GameState.Statistics.TreasureCollected))
	g.ui.AddLog("Starting new game...")
	g.startNewGame()
}

func (g *Game) SaveStatistics() error {
	if g.World == nil || g.World.GameState == nil {
		return nil
	}
	playthrough := g.World.GameState.ToPlaythroughStatistics()
	return g.statisticsStore.SavePlaythrough(playthrough)
}

func (g *Game) SaveGameState() error {
	if g.World == nil || g.World.GameState == nil {
		return fmt.Errorf("world or game state is nil")
	}
	err := g.gameStateStore.Save(g.World)
	if err == nil {
		g.ui.AddLog("Game saved successfully")
	}
	return err
}

func (g *Game) LoadGameState(playerName string) (*domain.World, error) {
	return g.gameStateStore.Load(playerName)
}

func (g *Game) HasSaveGame() bool {
	saves, err := g.gameStateStore.ListSaves()
	if err != nil {
		return false
	}
	return len(saves) > 0
}

func (g *Game) createNewGameWithPlayerName(playerName string) {
	if playerName == "" {
		playerName = "Player"
	}
	world := domain.NewWorldWithPlayerName(domain.Width, domain.Height, playerName)
	g.generator.Generate(world)
	g.World = world
	g.UpdateVisibility()
	g.ui.AddLog("New game started for " + playerName)
}

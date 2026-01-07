package service

import (
	"fmt"

	"rogue/internal/domain"
)

func (g *Game) loadGame() {
	saves, err := g.gameStateStore.ListSaves()
	if err != nil || len(saves) == 0 {
		g.ui.AddLog("No save")
		g.startNewGame()
		return
	}

	playerName := saves[0]
	world, err := g.LoadGameState(playerName)
	if err != nil || world == nil {
		g.ui.AddLog("Load failed")
		g.startNewGame()
		return
	}

	g.World = world
	g.UpdateVisibility()
	g.ui.AddLog("Loaded")
}

func (g *Game) startNewGame() {
	g.ui.SetNameInputState(true)
}

func (g *Game) handleGameCompletion() {
	g.SaveStatistics()
	g.ui.AddLog("Victory!")
	g.ui.AddLog(fmt.Sprintf("Score: $%d", g.World.GameState.Statistics.TreasureCollected))
	g.startNewGame()
}

func (g *Game) SaveStatistics() error {
	if g.World == nil || g.World.GameState == nil {
		return nil
	}

	playthrough := g.World.GameState.ToPlayStats(g.World.Player.Name)
	return g.statisticsStore.SavePlaythrough(playthrough)
}

func (g *Game) SaveGameState() error {
	if g.World == nil || g.World.GameState == nil {
		return fmt.Errorf("world or game state is nil")
	}
	err := g.gameStateStore.Save(g.World)
	if err == nil {
		g.ui.AddLog("Saved")
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
	world := domain.NewWorld(domain.WorldConfig.Width, domain.WorldConfig.Height, playerName)
	g.generator.Generate(world)
	g.World = world
	g.UpdateVisibility()
	g.ui.AddLog("New game")
}

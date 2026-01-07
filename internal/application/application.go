package application

import (
	"fmt"

	"rogue/internal/adapters/storage"
	"rogue/internal/adapters/tui"
	"rogue/internal/service"
)

type Application struct {
	gui             service.UI
	gameStateStore  *storage.StateStore
	statisticsStore *storage.StatsStore
}

func NewApplication() *Application {
	window := tui.NewWindow()
	gameStateStore := storage.NewStateStore("saves")
	statisticsStore := storage.NewStatsStore("saves")

	return &Application{
		gui:             window,
		gameStateStore:  gameStateStore,
		statisticsStore: statisticsStore,
	}
}

func (a *Application) Run() error {
	if err := a.gui.Init(); err != nil {
		return fmt.Errorf("failed to init gui: %w", err)
	}
	defer a.gui.Close()

	game := service.NewGame(a.gui, a.gameStateStore, a.statisticsStore)
	game.Run()

	return nil
}

package application

import (
	"fmt"

	"rogue/internal/adapters/storage"
	"rogue/internal/adapters/tui"
	"rogue/internal/service"
)

type Application struct {
	gui   service.UI
	state service.StateStore
	stats service.StatsStore
}

func NewApplication() *Application {
	window := tui.NewWindow()
	state := storage.NewStateStore("saves")
	stats := storage.NewStatsStore("saves")

	return &Application{
		gui:   window,
		state: state,
		stats: stats,
	}
}

func (a *Application) Run() error {
	if err := a.gui.Init(); err != nil {
		return fmt.Errorf("failed to init gui: %w", err)
	}
	defer a.gui.Close()

	game := service.NewGame(a.gui, a.state, a.stats)
	game.Run()

	return nil
}

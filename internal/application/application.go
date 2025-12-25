package application

import (
	"fmt"

	"rogue/internal/adapters/tui"
	"rogue/internal/service"
)

type Application struct {
	gui service.UI
}

func NewApplication() *Application {
	window := tui.NewWindow()

	return &Application{
		gui: window,
	}
}

func (a *Application) Run() error {
	if err := a.gui.Init(); err != nil {
		return fmt.Errorf("failed to init gui: %w", err)
	}
	defer a.gui.Close()

	game := service.NewGame(a.gui)
	game.Run()

	return nil
}

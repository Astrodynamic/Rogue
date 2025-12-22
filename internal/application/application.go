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

// Run - точка входа в процесс
func (a *Application) Run() error {
	// 1. Инициализация экрана
	if err := a.gui.Init(); err != nil {
		return fmt.Errorf("failed to init gui: %w", err)
	}
	// Гарантируем закрытие
	defer a.gui.Close()

	// 2. Игровой цикл (Game Loop)
	// Пока он живет прямо тут, позже переедет в service.Game
	running := true
	counter := 0

	for running {
		// A. Render
		a.gui.Draw(fmt.Sprintf("Frame: %d | Press Q to Quit", counter))

		// B. Input
		cmd := a.gui.PollInput()

		// C. Update
		switch cmd {
		case service.CmdQuit:
			running = false
		case service.CmdMove:
			counter++
		}
	}

	return nil
}

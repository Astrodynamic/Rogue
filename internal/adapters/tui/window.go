package tui

import (
	"rogue/internal/service"

	"github.com/gdamore/tcell/v2"
)

type Window struct {
	screen tcell.Screen
}

func NewWindow() *Window {
	return &Window{}
}

func (w *Window) Init() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := s.Init(); err != nil {
		return err
	}
	w.screen = s
	return nil
}

func (w *Window) Close() {
	if w.screen != nil {
		w.screen.Fini()
	}
}

func (w *Window) Draw(text string) {
	w.screen.Clear()

	// Простейший вывод текста для теста (в центре экрана)
	// (В реальности тут будет цикл по тайлам)
	for i, r := range text {
		w.screen.SetContent(i+1, 1, r, nil, tcell.StyleDefault)
	}

	w.screen.Show()
}

func (w *Window) PollInput() service.Command {
	ev := w.screen.PollEvent()

	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
			return service.CmdQuit
		}
		// Любая другая кнопка считаем за движение для теста
		return service.CmdMove
	}
	return service.CmdNone
}

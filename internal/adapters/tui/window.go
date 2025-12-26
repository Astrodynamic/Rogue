package tui

import (
	"rogue/internal/domain"
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

func (w *Window) Draw(world *domain.World) {
	w.screen.Clear()
	w.DrawWorld(world)
	w.screen.Show()
}

func (w *Window) Input() service.Command {
	ev := w.screen.PollEvent()

	switch ev := ev.(type) {
	case *tcell.EventKey:
		if cmd := w.KeyMap(ev.Key()); cmd != service.CmdNone {
			return cmd
		}

		if cmd := w.RuneMap(ev.Rune()); cmd != service.CmdNone {
			return cmd
		}
	}
	return service.CmdNone
}

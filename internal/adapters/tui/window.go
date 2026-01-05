package tui

import (
	"rogue/internal/domain"
	"rogue/internal/service"

	"github.com/gdamore/tcell/v2"
)

type Layout struct {
	Screen    domain.Rect
	Map       domain.Rect
	Stats     domain.Rect
	Inventory domain.Rect
	Log       domain.Rect
}

type Window struct {
	screen tcell.Screen
	log    *Log
	layout Layout
}

func NewWindow() *Window {
	return &Window{
		log: NewLog(10),
	}
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
	w.resize()
	return nil
}

func (w *Window) Close() {
	if w.screen != nil {
		w.screen.Fini()
	}
}

func (w *Window) Draw(world *domain.World) {
	w.screen.Clear()

	w.DrawMap(world, w.layout.Map)
	w.DrawStats(world, w.layout.Stats)
	w.DrawInventory(world.Player, w.layout.Inventory)
	w.DrawLog(w.layout.Log)

	w.screen.Show()
}

func (w *Window) resize() {
	wW, wH := w.screen.Size()
	if w.layout.Screen.W == wW && w.layout.Screen.H == wH {
		return
	}

	w.layout.Screen = domain.NewRect(0, 0, wW, wH)

	panelW := panelWidth(wW)
	logH := logHeight(wH)
	mapW := wW - panelW
	mapH := wH - logH
	statsH := 6
	invH := mapH - statsH

	w.layout.Map = domain.NewRect(0, 0, mapW, mapH)
	w.layout.Stats = domain.NewRect(mapW, 0, panelW, statsH)
	w.layout.Inventory = domain.NewRect(mapW, statsH, panelW, invH)
	w.layout.Log = domain.NewRect(0, mapH, wW, logH)
}

func panelWidth(width int) int {
	if width < 80 {
		return 20
	}
	return 28
}

func logHeight(height int) int {
	if height < 30 {
		return 5
	}
	return 8
}

func (w *Window) Input() service.Command {
	ev := w.screen.PollEvent()

	switch ev := ev.(type) {
	case *tcell.EventResize:
		w.resize()
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

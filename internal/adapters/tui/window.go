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
	Equipment domain.Rect
	Inventory domain.Rect
	Log       domain.Rect
}

type Window struct {
	screen         tcell.Screen
	log            *Log
	layout         Layout
	showStatistics bool
	statisticsData []*domain.PlaythroughStatistics
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

func (w *Window) Draw(world *domain.World, selection service.SelectionState) {
	w.screen.Clear()

	if w.showStatistics {
		w.drawStatisticsView(w.statisticsData, w.layout.Screen)
	} else {
		w.DrawMap(world, w.layout.Map)
		w.DrawStats(world, w.layout.Stats)
		w.DrawEquipment(world.Player, w.layout.Equipment, selection)
		w.DrawInventory(world.Player, w.layout.Inventory, selection)
		w.DrawLog(w.layout.Log)
	}

	w.screen.Show()
}

func (w *Window) DrawStatistics(playthroughs []*domain.PlaythroughStatistics) {
	w.showStatistics = true
	w.statisticsData = playthroughs
	w.screen.Clear()
	w.drawStatisticsView(playthroughs, w.layout.Screen)
	w.screen.Show()
}

func (w *Window) HideStatistics() {
	w.showStatistics = false
	w.statisticsData = nil
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
	mapH := wH
	statsH := 6
	equipmentH := 8
	invH := wH - statsH - equipmentH - logH

	w.layout.Map = domain.NewRect(0, 0, mapW, mapH)
	w.layout.Stats = domain.NewRect(mapW, 0, panelW, statsH)
	w.layout.Equipment = domain.NewRect(mapW, statsH, panelW, equipmentH)
	w.layout.Inventory = domain.NewRect(mapW, statsH+equipmentH, panelW, invH)
	w.layout.Log = domain.NewRect(mapW, statsH+equipmentH+invH, panelW, logH)
}

func panelWidth(width int) int {
	if width < 80 {
		return 20
	}
	return 28
}

func logHeight(height int) int {
	if height < 30 {
		return 10
	}
	return 15
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

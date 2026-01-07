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
	menuState      *MenuState
	statisticsData []*domain.PlayStats
}

func NewWindow() *Window {
	return &Window{
		log:       NewLog(50),
		menuState: NewMenuState(),
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
	if w.menuState.IsNameInput() {
		w.DrawNameInput()
	} else if world != nil {
		w.DrawMap(world, w.layout.Map)
		w.DrawStats(world, w.layout.Stats)
		w.DrawEquipment(world.Player, w.layout.Equipment, selection)
		w.DrawInventory(world.Player, w.layout.Inventory, selection)
		w.DrawLog(w.layout.Log)
	}
	w.screen.Show()
}

func (w *Window) DrawStatistics(playthroughs []*domain.PlayStats) {
	w.menuState.SetStatistics(true)
	w.statisticsData = playthroughs
	w.screen.Clear()
	if !w.menuState.IsNameInput() {
		w.drawStatisticsView(playthroughs, w.layout.Screen)
	}
	w.screen.Show()
}

func (w *Window) IsStartMenu() bool {
	return w.menuState.IsStartMenu()
}

func (w *Window) IsStatistics() bool {
	return w.menuState.IsStatistics()
}

func (w *Window) IsNameInput() bool {
	return w.menuState.IsNameInput()
}

func (w *Window) GetPlayerNameInput() string {
	return w.menuState.GetPlayerName()
}

func (w *Window) SetNameInputState(enabled bool) {
	w.menuState.SetNameInput(enabled)
}

func (w *Window) GetMenuOption() int {
	return w.menuState.GetMenuOption()
}

func (w *Window) HandleMenuNavigation(cmd service.Command, hasSave bool) {
	maxOptions := 3
	if hasSave {
		maxOptions = 4
	}

	switch cmd {
	case service.CmdSelectUp:
		w.menuState.MoveUp(maxOptions)
	case service.CmdSelectDown:
		w.menuState.MoveDown(maxOptions)
	case service.CmdBackToGame, service.CmdSelectCancel, service.CmdQuit:
		if w.menuState.IsStatistics() {
			w.menuState.Reset()
		}
	}
}

func (w *Window) ProcessMenuSelection(hasSave bool) service.Command {
	option := w.menuState.GetMenuOption()
	maxOptions := 3
	if hasSave {
		maxOptions = 4
	}

	if option == maxOptions-1 {
		return service.CmdQuit
	}

	if hasSave {
		switch option {
		case 0:
			w.menuState.SetStartMenu(false)
			return service.CmdNewGame
		case 1:
			w.menuState.SetStartMenu(false)
			return service.CmdLoadGame
		case 2:
			w.menuState.SetStartMenu(false)
			w.menuState.SetStatistics(true)
			return service.CmdShowStatistics
		}
	} else {
		switch option {
		case 0:
			w.menuState.SetStartMenu(false)
			return service.CmdNewGame
		case 1:
			w.menuState.SetStartMenu(false)
			w.menuState.SetStatistics(true)
			return service.CmdShowStatistics
		}
	}

	return service.CmdNone
}

func (w *Window) ReturnToMenu() {
	w.menuState.Reset()
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
	statsH := domain.UIStatsHeight
	equipmentH := domain.UIEquipmentHeight
	invH := wH - statsH - equipmentH - logH

	w.layout.Map = domain.NewRect(0, 0, mapW, mapH)
	w.layout.Stats = domain.NewRect(mapW, 0, panelW, statsH)
	w.layout.Equipment = domain.NewRect(mapW, statsH, panelW, equipmentH)
	w.layout.Inventory = domain.NewRect(mapW, statsH+equipmentH, panelW, invH)
	w.layout.Log = domain.NewRect(mapW, statsH+equipmentH+invH, panelW, logH)
}

func panelWidth(width int) int {
	if width < domain.UIMinWidthThreshold {
		return domain.UIMinPanelWidth
	}
	return domain.UIDefaultPanelWidth
}

func logHeight(height int) int {
	if height < domain.UIMinHeightThreshold {
		return domain.UIMinLogHeight
	}
	return domain.UIDefaultLogHeight
}

func (w *Window) Input() service.Command {
	ev := w.screen.PollEvent()

	switch ev := ev.(type) {
	case *tcell.EventResize:
		w.resize()
	case *tcell.EventKey:
		if w.menuState.IsNameInput() {
			cmd := w.KeyMap(ev.Key())
			var char rune = ev.Rune()
			if cmd != service.CmdNone && cmd != service.CmdBackToGame {
				char = 0
			}
			w.HandleNameInput(cmd, char)
			if cmd == service.CmdSelectConfirm {
				return service.CmdSelectConfirm
			}
			if cmd == service.CmdSelectCancel || cmd == service.CmdQuit {
				return cmd
			}
			return service.CmdNone
		}

		if cmd := w.KeyMap(ev.Key()); cmd != service.CmdNone {
			return cmd
		}

		if cmd := w.RuneMap(ev.Rune()); cmd != service.CmdNone {
			return cmd
		}
	}
	return service.CmdNone
}

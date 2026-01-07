package tui

import (
	"rogue/internal/service"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) KeyMap(key tcell.Key) service.Command {
	switch key {
	case tcell.KeyEscape:
		if w.menuState.IsStatistics() {
			return service.CmdBackToGame
		}
		return service.CmdSelectCancel
	case tcell.KeyUp:
		return service.CmdSelectUp
	case tcell.KeyDown:
		return service.CmdSelectDown
	case tcell.KeyEnter:
		return service.CmdSelectConfirm
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if w.menuState.IsNameInput() {
			return service.CmdBackToGame
		}
	}
	return service.CmdNone
}

func (w *Window) RuneMap(rune rune) service.Command {
	if w.menuState.IsStatistics() {
		return service.CmdNone
	}

	switch rune {
	case 'q', 'Q':
		return service.CmdQuit
	case 'w', 'W':
		return service.CmdMoveUp
	case 's', 'S':
		return service.CmdMoveDown
	case 'a', 'A':
		return service.CmdMoveLeft
	case 'd', 'D':
		return service.CmdMoveRight
	case 'h', 'H':
		return service.CmdUseWeapon
	case 'r', 'R':
		return service.CmdUseArmor
	case 'j', 'J':
		return service.CmdUseFood
	case 'k', 'K':
		return service.CmdUseElixir
	case 'e', 'E':
		return service.CmdUseScroll
	case 'i', 'I':
		return service.CmdEquipItem
	case 'u', 'U':
		return service.CmdUnequipItem
	case 'x', 'X':
		return service.CmdDropItem
	}
	return service.CmdNone
}

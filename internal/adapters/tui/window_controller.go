package tui

import (
	"rogue/internal/service"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) KeyMap(key tcell.Key) service.Command {
	switch key {
	case tcell.KeyEscape:
		return service.CmdQuit
	case tcell.KeyUp:
		return service.CmdMoveUp
	case tcell.KeyDown:
		return service.CmdMoveDown
	case tcell.KeyLeft:
		return service.CmdMoveLeft
	case tcell.KeyRight:
		return service.CmdMoveRight
	}
	return service.CmdNone
}

func (w *Window) RuneMap(rune rune) service.Command {
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
	}
	return service.CmdNone
}

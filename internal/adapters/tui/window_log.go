package tui

import (
	"rogue/internal/domain"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) DrawLog(rect domain.Rect) {
	w.drawBox(rect, "Log")

	messages := w.log.GetMessages()
	startLine := 0
	if len(messages) > rect.H-2 {
		startLine = len(messages) - (rect.H - 2)
	}

	line := 1
	style := tcell.StyleDefault.Foreground(tcell.ColorGray)
	for i := startLine; i < len(messages) && line < rect.H-1; i++ {
		msg := messages[i]
		if len(msg) > rect.W-2 {
			msg = msg[:rect.W-2]
		}
		for j, r := range msg {
			if j >= rect.W-2 {
				break
			}
			w.screen.SetContent(rect.X+1+j, rect.Y+line, r, nil, style)
		}
		line++
	}
}

func (w *Window) AddLog(message string) {
	w.log.Add(message)
}

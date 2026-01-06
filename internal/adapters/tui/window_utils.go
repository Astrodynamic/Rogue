package tui

import (
	"fmt"
	"rogue/internal/domain"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) drawBox(rect domain.Rect, title string) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	for i := 0; i < rect.W; i++ {
		w.screen.SetContent(rect.X+i, rect.Y, '─', nil, style)
		w.screen.SetContent(rect.X+i, rect.Y+rect.H-1, '─', nil, style)
	}

	for i := 0; i < rect.H; i++ {
		w.screen.SetContent(rect.X, rect.Y+i, '│', nil, style)
		w.screen.SetContent(rect.X+rect.W-1, rect.Y+i, '│', nil, style)
	}

	w.screen.SetContent(rect.X, rect.Y, '┌', nil, style)
	w.screen.SetContent(rect.X+rect.W-1, rect.Y, '┐', nil, style)
	w.screen.SetContent(rect.X, rect.Y+rect.H-1, '└', nil, style)
	w.screen.SetContent(rect.X+rect.W-1, rect.Y+rect.H-1, '┘', nil, style)

	if title != "" {
		titleText := " " + title + " "
		titleX := rect.X + (rect.W-len(titleText))/2
		for i, r := range titleText {
			if titleX+i < rect.X+rect.W-1 {
				w.screen.SetContent(titleX+i, rect.Y, r, nil, style)
			}
		}
	}
}

func (w *Window) drawText(rect domain.Rect, line int, format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	w.drawTextRaw(rect, line, text, tcell.StyleDefault.Foreground(tcell.ColorWhite))
}

func (w *Window) drawTextHighlighted(rect domain.Rect, line int, text string) {
	w.drawTextRaw(rect, line, text, tcell.StyleDefault.Reverse(true))
}

func (w *Window) drawTextRaw(rect domain.Rect, line int, text string, style tcell.Style) {
	width := rect.W - 2
	if len(text) > width {
		text = text[:width]
	}

	x := rect.X + 1
	y := rect.Y + line
	for i, r := range text {
		if i >= width {
			break
		}
		w.screen.SetContent(x+i, y, r, nil, style)
	}
}

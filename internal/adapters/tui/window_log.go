package tui

import (
	"rogue/internal/domain"

	"github.com/gdamore/tcell/v2"
)

type Log struct {
	messages []string
	maxLines int
}

func NewLog(maxLines int) *Log {
	return &Log{
		messages: make([]string, 0, maxLines*2),
		maxLines: maxLines * 2,
	}
}

func (l *Log) Add(message string) {
	if message == "" {
		return
	}

	l.messages = append(l.messages, message)
	if len(l.messages) > l.maxLines {
		l.messages = l.messages[1:]
	}
}

func (l *Log) GetMessages() []string {
	return l.messages
}

func (l *Log) Clear() {
	l.messages = l.messages[:0]
}

func (w *Window) DrawLog(rect domain.Rect) {
	w.drawBox(rect, "Log")

	messages := w.log.GetMessages()
	availableLines := rect.H - 2
	if availableLines < 1 {
		return
	}

	startLine := 0
	if len(messages) > availableLines {
		startLine = len(messages) - availableLines
	}

	line := 1
	style := tcell.StyleDefault.Foreground(tcell.ColorGray)
	maxWidth := rect.W - 2

	for i := startLine; i < len(messages) && line < rect.H-1; i++ {
		msg := messages[i]
		if len(msg) > maxWidth {
			msg = msg[:maxWidth]
		}

		x := rect.X + 1
		for _, r := range msg {
			if x >= rect.X+rect.W-1 {
				break
			}
			w.screen.SetContent(x, rect.Y+line, r, nil, style)
			x++
		}
		line++
	}
}

func (w *Window) AddLog(message string) {
	w.log.Add(message)
}

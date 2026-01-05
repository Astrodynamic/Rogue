package tui

type Log struct {
	messages []string
	maxLines int
}

func NewLog(maxLines int) *Log {
	return &Log{
		messages: make([]string, 0, maxLines),
		maxLines: maxLines,
	}
}

func (l *Log) Add(message string) {
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

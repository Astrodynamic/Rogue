package service

type Command int

const (
	CmdNone Command = iota
	CmdQuit
	CmdMove
)

type UI interface {
	Init() error
	Close()

	Draw(text string)
	PollInput() Command
}

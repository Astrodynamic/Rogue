package service

import "rogue/internal/domain"

type Command int

const (
	CmdNone Command = iota
	CmdQuit
	CmdMoveUp
	CmdMoveDown
	CmdMoveLeft
	CmdMoveRight
)

type UI interface {
	Init() error
	Close()

	Draw(world *domain.World)
	Input() Command
}

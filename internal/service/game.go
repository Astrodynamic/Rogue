package service

import (
	"rogue/internal/domain"
)

type Game struct {
	World     *domain.World
	isRunning bool
	ui        UI
}

func NewGame(ui UI) *Game {
	h, w := 25, 80

	return &Game{
		World:     domain.NewWorld(h, w),
		isRunning: true,
		ui:        ui,
	}
}

func (g *Game) Run() {
	for g.isRunning {
		g.ui.Draw(g.World)
		g.handle(g.ui.Input())
	}
}

func (g *Game) handle(cmd Command) {
	switch cmd {
	case CmdQuit:
		g.isRunning = false
	case CmdMoveUp:
		g.onPlayerMove(domain.DirUp)
	case CmdMoveDown:
		g.onPlayerMove(domain.DirDown)
	case CmdMoveLeft:
		g.onPlayerMove(domain.DirLeft)
	case CmdMoveRight:
		g.onPlayerMove(domain.DirRight)
	}
}

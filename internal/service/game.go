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

	world := domain.NewWorld(w, h)

	NewGenerator().Generate(world)

	return &Game{
		World:     world,
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
		g.onMove(&g.World.Player.Actor, domain.DirUp)
	case CmdMoveDown:
		g.onMove(&g.World.Player.Actor, domain.DirDown)
	case CmdMoveLeft:
		g.onMove(&g.World.Player.Actor, domain.DirLeft)
	case CmdMoveRight:
		g.onMove(&g.World.Player.Actor, domain.DirRight)
	}
}

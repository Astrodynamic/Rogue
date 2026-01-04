package service

import (
	"rogue/internal/domain"
)

type Game struct {
	World     *domain.World
	FOV       *FOV
	isRunning bool
	ui        UI
}

func NewGame(ui UI) *Game {
	world := domain.NewWorld(domain.Width, domain.Height)

	NewGenerator().Generate(world)

	game := &Game{
		World:     world,
		FOV:       NewFOV(),
		isRunning: true,
		ui:        ui,
	}

	game.UpdateVisibility()

	return game
}

func (g *Game) UpdateVisibility() {
	g.FOV.Update(g.World.Level, g.World.Player.Point)
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
		g.onMove(&g.World.Player.Actor, domain.DirSU)
	case CmdMoveDown:
		g.onMove(&g.World.Player.Actor, domain.DirSD)
	case CmdMoveLeft:
		g.onMove(&g.World.Player.Actor, domain.DirLS)
	case CmdMoveRight:
		g.onMove(&g.World.Player.Actor, domain.DirRS)
	}
}

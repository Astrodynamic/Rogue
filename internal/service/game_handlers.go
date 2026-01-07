package service

import (
	"rogue/internal/domain"
)

func (g *Game) handle(cmd Command) {
	if g.ui.IsStartMenu() {
		g.handleStartMenu(cmd)
		return
	}

	if g.ui.IsStatistics() {
		g.handleStatisticsView(cmd)
		return
	}

	if g.selection != nil && g.selection.IsActive() {
		g.handleSelection(cmd)
		return
	}

	g.handleGameCommand(cmd)
}

func (g *Game) handleStartMenu(cmd Command) {
	g.ui.HandleMenuNavigation(cmd, g.HasSaveGame())

	if cmd == CmdSelectConfirm {
		menuCmd := g.ui.ProcessMenuSelection(g.HasSaveGame())
		switch menuCmd {
		case CmdNewGame:
			g.startNewGame()
		case CmdLoadGame:
			g.loadGame()
		case CmdShowStatistics:
		case CmdQuit:
			g.isRunning = false
		}
	}
}

func (g *Game) handleStatisticsView(cmd Command) {
	g.ui.HandleMenuNavigation(cmd, g.HasSaveGame())
}

func (g *Game) handleGameCommand(cmd Command) {
	switch cmd {
	case CmdQuit:
		g.returnToMenu()
	case CmdMoveUp:
		g.onMove(&g.World.Player.Actor, domain.DirSU)
	case CmdMoveDown:
		g.onMove(&g.World.Player.Actor, domain.DirSD)
	case CmdMoveLeft:
		g.onMove(&g.World.Player.Actor, domain.DirLS)
	case CmdMoveRight:
		g.onMove(&g.World.Player.Actor, domain.DirRS)
	case CmdUseWeapon:
		g.startItemSelection(domain.ItemWeapon, false)
	case CmdUseArmor:
		g.startItemSelection(domain.ItemArmor, false)
	case CmdUseFood:
		g.startItemSelection(domain.ItemFood, false)
	case CmdUseElixir:
		g.startItemSelection(domain.ItemElixir, false)
	case CmdUseScroll:
		g.startItemSelection(domain.ItemScroll, false)
	case CmdEquipItem:
		g.startEquipItemSelection()
	case CmdUnequipItem:
		g.startUnequipItemSelection()
	case CmdDropItem:
		g.startDropItemSelection()
	}
}

func (g *Game) returnToMenu() {
	g.SaveStatistics()
	if g.World != nil {
		g.SaveGameState()
	}
	g.ui.ReturnToMenu()
}

func (g *Game) handleSelection(cmd Command) {
	if g.selection == nil {
		return
	}

	switch cmd {
	case CmdSelectCancel:
		g.selection.Cancel()
		g.selection = nil
	case CmdSelectUp:
		g.selection.MoveUp()
	case CmdSelectDown:
		g.selection.MoveDown()
	case CmdSelectConfirm:
		if g.selection.Confirm() {
			g.selection = nil
		}
	}
}

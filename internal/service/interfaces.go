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
	CmdUseWeapon
	CmdUseArmor
	CmdUseFood
	CmdUseElixir
	CmdUseScroll
	CmdEquipItem
	CmdUnequipItem
	CmdDropItem
	CmdSelectUp
	CmdSelectDown
	CmdSelectConfirm
	CmdSelectCancel
	CmdShowStatistics
	CmdBackToGame
)

type SelectionState struct {
	Model *SelectionModel
}

type UI interface {
	Init() error
	Close()

	Draw(world *domain.World, selection SelectionState)
	DrawStatistics(playthroughs []*domain.PlaythroughStatistics)
	HideStatistics()
	Input() Command
}

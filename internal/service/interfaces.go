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
	CmdBackToGame
	CmdNewGame
	CmdLoadGame
	CmdShowStatistics
)

type SelectionState struct {
	Model *SelectionModel
}

type UI interface {
	Init() error
	Close()
	Draw(world *domain.World, selection SelectionState)
	DrawStatistics(playthroughs []*domain.PlaythroughStatistics)
	DrawStartMenu(hasSave bool, currentOption int)
	AddLog(message string)
	Input() Command
	IsStartMenu() bool
	IsStatistics() bool
	IsNameInput() bool
	GetPlayerNameInput() string
	SetNameInputState(enabled bool)
	GetMenuOption() int
	HandleMenuNavigation(cmd Command, hasSave bool)
	ProcessMenuSelection(hasSave bool) Command
	ReturnToMenu()
}

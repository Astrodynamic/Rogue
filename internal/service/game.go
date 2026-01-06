package service

import (
	"rogue/internal/domain"
)

type Game struct {
	World     *domain.World
	FOV       *FOV
	isRunning bool
	ui        UI
	selection *SelectionModel
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
		selection := SelectionState{
			Model: g.selection,
		}
		g.ui.Draw(g.World, selection)
		g.handle(g.ui.Input())
	}
}

func (g *Game) handle(cmd Command) {
	if g.selection != nil && g.selection.IsActive() {
		g.handleSelection(cmd)
		return
	}

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
	case CmdSelectUp, CmdSelectDown, CmdSelectConfirm, CmdSelectCancel:
	case CmdUseWeapon:
		g.startItemSelection(domain.ItemWeapon, false)
	case CmdUseFood:
		g.startItemSelection(domain.ItemFood, false)
	case CmdUseElixir:
		g.startItemSelection(domain.ItemElixir, false)
	case CmdUseScroll:
		g.startItemSelection(domain.ItemScroll, false)
	case CmdDropItem:
		g.startDropItemSelection()
	case CmdDropEquipment:
		g.startDropEquipmentSelection()
	}
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

func (g *Game) startItemSelection(kind domain.ItemKind, isDrop bool) {
	stacks := g.World.Player.Backpack.GetStacks(kind)
	if len(stacks) == 0 {
		return
	}

	selectionItems := make([]SelectionItem, 0)
	for i, stack := range stacks {
		if item := stack.Peek(); item != nil {
			selectionItems = append(selectionItems, &ItemSelectionItem{
				Item:       item,
				ItemKind:   kind,
				StackIndex: i,
			})
		}
	}

	if len(selectionItems) == 0 {
		return
	}

	var handler SelectionHandler
	if isDrop {
		handler = NewDropItemHandler(g)
	} else {
		handler = NewUseItemHandler(g, kind)
	}

	g.selection = NewSelectionModel(selectionItems, handler)
}

func (g *Game) startDropItemSelection() {
	allKinds := []domain.ItemKind{
		domain.ItemFood,
		domain.ItemElixir,
		domain.ItemScroll,
		domain.ItemWeapon,
		domain.ItemArmor,
	}

	selectionItems := make([]SelectionItem, 0)

	for _, kind := range allKinds {
		stacks := g.World.Player.Backpack.GetStacks(kind)
		for stackIndex, stack := range stacks {
			if item := stack.Peek(); item != nil {
				selectionItems = append(selectionItems, &ItemSelectionItem{
					Item:       item,
					ItemKind:   kind,
					StackIndex: stackIndex,
				})
			}
		}
	}

	if len(selectionItems) == 0 {
		return
	}

	g.selection = NewSelectionModel(selectionItems, NewDropItemHandler(g))
}

func (g *Game) startDropEquipmentSelection() {
	parts := []domain.ActorPart{
		domain.ActorPartHead,
		domain.ActorPartBody,
		domain.ActorPartHand,
		domain.ActorPartLegs,
	}

	selectionItems := make([]SelectionItem, 0)
	for _, part := range parts {
		item := g.World.Player.Equipment.Get(part)
		if item != nil {
			selectionItems = append(selectionItems, &EquipmentSelectionItem{
				Part: part,
				Item: item,
			})
		}
	}

	if len(selectionItems) == 0 {
		return
	}

	g.selection = NewSelectionModel(selectionItems, NewUnequipEquipmentHandler(g))
}

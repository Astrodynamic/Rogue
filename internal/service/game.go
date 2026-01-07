package service

import (
	"math/rand/v2"
	"time"

	"rogue/internal/adapters/storage"
	"rogue/internal/domain"
)

type Game struct {
	World           *domain.World
	FOV             *FOV
	isRunning       bool
	ui              UI
	selection       *SelectionModel
	statisticsStore *storage.StatisticsStorage
	showStatistics  bool
	rng             domain.RandomGenerator
}

func NewGame(ui UI) *Game {
	world := domain.NewWorld(domain.Width, domain.Height)

	NewGenerator().Generate(world)

	// Create random generator for combat
	seed := uint64(time.Now().UnixNano())
	randSource := rand.New(rand.NewPCG(seed, seed))
	rng := NewRandAdapter(randSource)

	game := &Game{
		World:           world,
		FOV:             NewFOV(),
		isRunning:       true,
		ui:              ui,
		statisticsStore: storage.NewStatisticsStorage("saves"),
		showStatistics:  false,
		rng:             rng,
	}

	game.UpdateVisibility()

	return game
}

func (g *Game) UpdateVisibility() {
	g.FOV.Update(g.World.Level, g.World.Player.Point)
}

func (g *Game) Run() {
	for g.isRunning {
		if g.showStatistics {
			playthroughs, _ := g.statisticsStore.GetLeaderboard(0)
			g.ui.DrawStatistics(playthroughs)
		} else {
			selection := SelectionState{
				Model: g.selection,
			}
			g.ui.Draw(g.World, selection)
		}
		g.handle(g.ui.Input())
	}
}

func (g *Game) handle(cmd Command) {
	if g.showStatistics {
		switch cmd {
		case CmdBackToGame, CmdSelectCancel:
			g.showStatistics = false
			g.ui.HideStatistics()
			selection := SelectionState{
				Model: g.selection,
			}
			g.ui.Draw(g.World, selection)
		}
		return
	}

	if g.selection != nil && g.selection.IsActive() {
		g.handleSelection(cmd)
		return
	}

	switch cmd {
	case CmdQuit:
		g.SaveStatistics()
		g.isRunning = false
	case CmdShowStatistics:
		g.showStatistics = true
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

func (g *Game) startEquipItemSelection() {
	equippableKinds := []domain.ItemKind{
		domain.ItemWeapon,
		domain.ItemArmor,
	}

	selectionItems := make([]SelectionItem, 0)
	for _, kind := range equippableKinds {
		stacks := g.World.Player.Backpack.GetStacks(kind)
		for stackIndex, stack := range stacks {
			if item := stack.Peek(); item != nil {
				_, canEquip := domain.GetEquipPart(item)
				if canEquip {
					selectionItems = append(selectionItems, &ItemSelectionItem{
						Item:       item,
						ItemKind:   kind,
						StackIndex: stackIndex,
					})
				}
			}
		}
	}

	if len(selectionItems) == 0 {
		return
	}

	g.selection = NewSelectionModel(selectionItems, NewEquipItemHandler(g))
}

func (g *Game) startUnequipItemSelection() {
	parts := domain.GetAllParts()

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

func (g *Game) RecordTreasureCollected(amount int) {
	g.World.GameState.RecordTreasureCollected(amount)
}

func (g *Game) RecordEnemyDefeated() {
	g.World.GameState.RecordEnemyDefeated()
}

func (g *Game) RecordFoodConsumed() {
	g.World.GameState.RecordFoodConsumed()
}

func (g *Game) RecordElixirDrunk() {
	g.World.GameState.RecordElixirDrunk()
}

func (g *Game) RecordScrollRead() {
	g.World.GameState.RecordScrollRead()
}

func (g *Game) RecordHitDealt() {
	g.World.GameState.RecordHitDealt()
}

func (g *Game) RecordHitReceived() {
	g.World.GameState.RecordHitReceived()
}

func (g *Game) RecordTileTraveled() {
	g.World.GameState.RecordTileTraveled()
}

func (g *Game) SaveStatistics() error {
	playthrough := g.World.GameState.ToPlaythroughStatistics()
	return g.statisticsStore.SavePlaythrough(playthrough)
}

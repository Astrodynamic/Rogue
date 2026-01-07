package service

import (
	"rogue/internal/domain"
)

func (g *Game) startItemSelection(kind domain.ItemKind, isDrop bool) {
	stacks := g.World.Player.Backpack.GetStacks(kind)
	if len(stacks) == 0 {
		return
	}

	items := g.buildSelectionItemsFromStacks(stacks, kind)
	if len(items) == 0 {
		return
	}

	var handler SelectionHandler
	if isDrop {
		handler = NewDropItemHandler(g)
	} else {
		handler = NewUseItemHandler(g, kind)
	}

	g.sel = NewSelectionModel(items, handler)
}

func (g *Game) startDropItemSelection() {
	allKinds := []domain.ItemKind{
		domain.ItemFood,
		domain.ItemElixir,
		domain.ItemScroll,
		domain.ItemWeapon,
		domain.ItemArmor,
	}

	items := make([]SelectionItem, 0)
	for _, kind := range allKinds {
		stacks := g.World.Player.Backpack.GetStacks(kind)
		for stackIndex, stack := range stacks {
			if item := stack.Peek(); item != nil {
				items = append(items, &ItemSelItem{
					Item:       item,
					ItemKind:   kind,
					StackIndex: stackIndex,
				})
			}
		}
	}

	if len(items) == 0 {
		return
	}

	g.sel = NewSelectionModel(items, NewDropItemHandler(g))
}

func (g *Game) startEquipItemSelection() {
	equippableKinds := []domain.ItemKind{
		domain.ItemWeapon,
		domain.ItemArmor,
	}

	items := make([]SelectionItem, 0)
	for _, kind := range equippableKinds {
		stacks := g.World.Player.Backpack.GetStacks(kind)
		for stackIndex, stack := range stacks {
			if item := stack.Peek(); item != nil {
				_, canEquip := domain.GetEquipPart(item)
				if canEquip {
					items = append(items, &ItemSelItem{
						Item:       item,
						ItemKind:   kind,
						StackIndex: stackIndex,
					})
				}
			}
		}
	}

	if len(items) == 0 {
		return
	}

		g.sel = NewSelectionModel(items, NewEquipHdl(g))
}

func (g *Game) startUnequipItemSelection() {
	parts := domain.GetAllParts()

	items := make([]SelectionItem, 0)
	for _, part := range parts {
		item := g.World.Player.Equipment.Get(part)
		if item != nil {
			items = append(items, &EquipSelItem{
				Part: part,
				Item: item,
			})
		}
	}

	if len(items) == 0 {
		return
	}

		g.sel = NewSelectionModel(items, NewUnequipHdl(g))
}

func (g *Game) buildSelectionItemsFromStacks(stacks []*domain.ItemStack, kind domain.ItemKind) []SelectionItem {
	items := make([]SelectionItem, 0)
	for i, stack := range stacks {
		if item := stack.Peek(); item != nil {
			items = append(items, &ItemSelItem{
				Item:       item,
				ItemKind:   kind,
				StackIndex: i,
			})
		}
	}
	return items
}

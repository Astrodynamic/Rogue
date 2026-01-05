package tui

import (
	"rogue/internal/domain"
)

func (w *Window) DrawInventory(player *domain.Player, rect domain.Rect) {
	w.drawBox(rect, "Inventory")

	line := 1
	itemKinds := []domain.ItemKind{
		domain.ItemFood,
		domain.ItemElixir,
		domain.ItemScroll,
		domain.ItemWeapon,
		domain.ItemTreasure,
	}

	for _, kind := range itemKinds {
		stacks := player.Backpack.GetStacks(kind)
		if len(stacks) == 0 {
			continue
		}

		stack := stacks[0]
		item := stack.Peek()
		if item == nil {
			continue
		}

		itemName := item.Name()
		var value int
		if kind == domain.ItemTreasure {
			value = player.Backpack.GetTotalTreasure()
		} else {
			value = player.Backpack.Count(kind)
		}
		w.drawText(rect, line, "%s: %d", itemName, value)
		line++

		if line >= rect.H-2 {
			break
		}
	}
}

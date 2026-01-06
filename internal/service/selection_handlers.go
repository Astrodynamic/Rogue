package service

import "rogue/internal/domain"

type UseItemHandler struct {
	game     *Game
	itemKind domain.ItemKind
}

func NewUseItemHandler(game *Game, itemKind domain.ItemKind) *UseItemHandler {
	return &UseItemHandler{
		game:     game,
		itemKind: itemKind,
	}
}

func (h *UseItemHandler) OnConfirm(item SelectionItem) bool {
	itemInfo, ok := item.(*ItemSelectionItem)
	if !ok {
		return false
	}

	stacks := h.game.World.Player.Backpack.GetStacks(h.itemKind)
	if itemInfo.StackIndex >= len(stacks) {
		return false
	}

	itemObj := itemInfo.Item
	if h.itemKind == domain.ItemWeapon {
		part, canEquip := domain.GetEquipPart(itemObj)
		if canEquip {
			oldItem := h.game.World.Player.EquipItem(part, itemObj)
			if oldItem != nil {
				dropPos := h.game.World.Level.GetAdjacentDropPoint(h.game.World.Player.Point)
				if dropPos.X >= 0 {
					h.game.World.Level.AddItem(dropPos, oldItem)
				}
			}
			h.game.World.Player.Backpack.Remove(h.itemKind, itemInfo.StackIndex)
		} else {
			result := h.game.World.Player.Backpack.Use(h.itemKind, itemInfo.StackIndex, &h.game.World.Player.Actor)
			if !result.Success {
				return false
			}
		}
	} else {
		result := h.game.World.Player.Backpack.Use(h.itemKind, itemInfo.StackIndex, &h.game.World.Player.Actor)
		if !result.Success {
			return false
		}
	}

	h.game.UpdateVisibility()
	return true
}

func (h *UseItemHandler) OnCancel() {
}

type DropItemHandler struct {
	game *Game
}

func NewDropItemHandler(game *Game) *DropItemHandler {
	return &DropItemHandler{game: game}
}

func (h *DropItemHandler) OnConfirm(item SelectionItem) bool {
	itemInfo, ok := item.(*ItemSelectionItem)
	if !ok {
		return false
	}

	droppedItem := h.game.World.Player.DropItem(itemInfo.ItemKind, itemInfo.StackIndex)
	if droppedItem != nil {
		dropPos := h.game.World.Level.GetAdjacentDropPoint(h.game.World.Player.Point)
		if dropPos.X >= 0 {
			h.game.World.Level.AddItem(dropPos, droppedItem)
		}
	}

	h.game.UpdateVisibility()
	return true
}

func (h *DropItemHandler) OnCancel() {
}

type UnequipEquipmentHandler struct {
	game *Game
}

func NewUnequipEquipmentHandler(game *Game) *UnequipEquipmentHandler {
	return &UnequipEquipmentHandler{game: game}
}

func (h *UnequipEquipmentHandler) OnConfirm(item SelectionItem) bool {
	equipInfo, ok := item.(*EquipmentSelectionItem)
	if !ok {
		return false
	}

	unequippedItem := h.game.World.Player.DropEquipment(equipInfo.Part)
	if unequippedItem != nil {
		if !h.game.World.Player.Backpack.Add(unequippedItem) {
			dropPos := h.game.World.Level.GetAdjacentDropPoint(h.game.World.Player.Point)
			if dropPos.X >= 0 {
				h.game.World.Level.AddItem(dropPos, unequippedItem)
			}
		}
	}

	h.game.UpdateVisibility()
	return true
}

func (h *UnequipEquipmentHandler) OnCancel() {
}

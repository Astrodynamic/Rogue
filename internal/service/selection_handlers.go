package service

import "rogue/internal/domain"

// equipItemFromBackpack equips an item from backpack and handles the old item
func equipItemFromBackpack(game *Game, part domain.ActorPart, itemObj domain.Item, itemKind domain.ItemKind, stackIndex int) bool {
	oldItem := game.World.Player.EquipItem(part, itemObj)
	if !game.World.Player.HandleOldItemOnEquip(oldItem, game.World.Level, game.World.Player.Point) {
		return false
	}
	game.World.Player.Backpack.Remove(itemKind, stackIndex)
	return true
}

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
	part, canEquip := domain.GetEquipPart(itemObj)
	if canEquip {
		if !equipItemFromBackpack(h.game, part, itemObj, h.itemKind, itemInfo.StackIndex) {
			return false
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
	if !h.game.World.Player.HandleOldItemOnEquip(unequippedItem, h.game.World.Level, h.game.World.Player.Point) {
		return false
	}

	h.game.UpdateVisibility()
	return true
}

func (h *UnequipEquipmentHandler) OnCancel() {
}

type EquipItemHandler struct {
	game *Game
}

func NewEquipItemHandler(game *Game) *EquipItemHandler {
	return &EquipItemHandler{game: game}
}

func (h *EquipItemHandler) OnConfirm(item SelectionItem) bool {
	itemInfo, ok := item.(*ItemSelectionItem)
	if !ok {
		return false
	}

	stacks := h.game.World.Player.Backpack.GetStacks(itemInfo.ItemKind)
	if itemInfo.StackIndex >= len(stacks) {
		return false
	}

	itemObj := itemInfo.Item
	part, canEquip := domain.GetEquipPart(itemObj)
	if !canEquip {
		return false
	}

	if !equipItemFromBackpack(h.game, part, itemObj, itemInfo.ItemKind, itemInfo.StackIndex) {
		return false
	}

	h.game.UpdateVisibility()
	return true
}

func (h *EquipItemHandler) OnCancel() {
}

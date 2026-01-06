package tui

import (
	"fmt"
	"rogue/internal/domain"
	"rogue/internal/service"
)

func (w *Window) DrawEquipment(player *domain.Player, rect domain.Rect, selection service.SelectionState) {
	w.drawBox(rect, "Equipment")

	line := 1
	parts := []domain.ActorPart{
		domain.ActorPartHead,
		domain.ActorPartBody,
		domain.ActorPartHand,
		domain.ActorPartLegs,
	}

	partNames := map[domain.ActorPart]string{
		domain.ActorPartHead: "Head",
		domain.ActorPartBody: "Body",
		domain.ActorPartHand: "Hand",
		domain.ActorPartLegs: "Legs",
	}

	equippedIndex := 0
	for _, part := range parts {
		if line >= rect.H-2 {
			return
		}

		item := player.Equipment.Get(part)
		partName := partNames[part]

		isSelected := false
		if selection.Model != nil && selection.Model.IsActive() {
			selectedItem := selection.Model.SelectedItem()
			if equipInfo, ok := selectedItem.(*service.EquipmentSelectionItem); ok {
				isSelected = equipInfo.Part == part && selection.Model.SelectedIndex() == equippedIndex
			}
		}

		var text string
		if item == nil {
			text = fmt.Sprintf("%s: -", partName)
		} else {
			itemName := item.Name()
			properties := w.formatItemProperties(item)
			if properties != "" {
				text = fmt.Sprintf("%s: %s %s", partName, itemName, properties)
			} else {
				text = fmt.Sprintf("%s: %s", partName, itemName)
			}
		}

		if isSelected {
			w.drawTextHighlighted(rect, line, text)
		} else {
			w.drawText(rect, line, text)
		}
		line++
		if item != nil {
			equippedIndex++
		}
	}
}

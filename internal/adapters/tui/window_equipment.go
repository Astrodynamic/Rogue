package tui

import (
	"rogue/internal/domain"
)

func (w *Window) DrawEquipment(player *domain.Player, rect domain.Rect) {
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

	for _, part := range parts {
		if line >= rect.H-2 {
			return
		}

		item := player.Equipment.Get(part)
		partName := partNames[part]

		if item == nil {
			w.drawText(rect, line, "%s: -", partName)
		} else {
			itemName := item.Name()
			properties := w.formatItemProperties(item)
			if properties != "" {
				w.drawText(rect, line, "%s: %s %s", partName, itemName, properties)
			} else {
				w.drawText(rect, line, "%s: %s", partName, itemName)
			}
		}
		line++
	}
}

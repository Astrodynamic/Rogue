package tui

import (
	"fmt"
	"rogue/internal/domain"
	"rogue/internal/service"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) DrawEquipment(player *domain.Player, rect domain.Rect, selection service.SelectionState) {
	w.drawBox(rect, "Equipment")

	line := 1
	parts := domain.GetAllParts()

	selectionPartMap := make(map[domain.ActorPart]int)
	if selection.Model != nil && selection.Model.IsActive() {
		items := selection.Model.Items()
		for i, selItem := range items {
			if equipInfo, ok := selItem.(*service.EquipmentSelectionItem); ok {
				selectionPartMap[equipInfo.Part] = i
			}
		}
	}

	for _, part := range parts {
		if line >= rect.H-2 {
			return
		}

		item := player.Equipment.Get(part)
		partName := domain.GetPartName(part)

		isSelected := false
		if selection.Model != nil && selection.Model.IsActive() && item != nil {
			if selIndex, exists := selectionPartMap[part]; exists {
				isSelected = selection.Model.SelectedIndex() == selIndex
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
			w.drawTextRaw(rect, line, text, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
		line++
	}
}

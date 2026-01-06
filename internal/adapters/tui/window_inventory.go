package tui

import (
	"fmt"
	"rogue/internal/domain"
	"rogue/internal/service"
)

func (w *Window) formatItemProperties(item domain.Item) string {
	stats := item.GetStats()
	itemType := item.Type()

	switch itemType {
	case domain.ItemFood:
		if stats.Health > 0 {
			return fmt.Sprintf("HP+%d", stats.Health)
		}
	case domain.ItemElixir:
		if elixir, ok := item.(*domain.Elixir); ok {
			var statName string
			switch elixir.ElixirKind {
			case domain.ElixirHealth:
				statName = "HP"
			case domain.ElixirMaxHealth:
				statName = "H"
			case domain.ElixirDexterity:
				statName = "D"
			case domain.ElixirStrength:
				statName = "S"
			default:
				statName = "?"
			}
			return fmt.Sprintf("%s+%d (%dt)", statName, elixir.Amount, elixir.Duration)
		}
	case domain.ItemScroll:
		if scroll, ok := item.(*domain.Scroll); ok {
			var statName string
			switch scroll.ScrollKind {
			case domain.ScrollHealth:
				statName = "HP"
			case domain.ScrollMaxHealth:
				statName = "H"
			case domain.ScrollDexterity:
				statName = "D"
			case domain.ScrollStrength:
				statName = "S"
			default:
				statName = "?"
			}
			return fmt.Sprintf("%s+%d", statName, scroll.Amount)
		}
	case domain.ItemWeapon:
		if stats.Strength > 0 {
			return fmt.Sprintf("S+%d", stats.Strength)
		}
	case domain.ItemArmor:
		var props []string
		if stats.Strength > 0 {
			props = append(props, fmt.Sprintf("S+%d", stats.Strength))
		}
		if stats.Dexterity > 0 {
			props = append(props, fmt.Sprintf("D+%d", stats.Dexterity))
		}
		if stats.MaxHealth > 0 {
			props = append(props, fmt.Sprintf("H+%d", stats.MaxHealth))
		}
		if len(props) == 0 {
			return "-"
		}
		result := props[0]
		for i := 1; i < len(props); i++ {
			result += " " + props[i]
		}
		return result
	case domain.ItemTreasure:
		if treasure, ok := item.(*domain.Treasure); ok {
			return fmt.Sprintf("$%d", treasure.Value)
		}
	}
	return ""
}

func (w *Window) DrawInventory(player *domain.Player, rect domain.Rect, selection service.SelectionState) {
	w.drawBox(rect, "Inventory")

	line := 1
	itemKinds := []domain.ItemKind{
		domain.ItemFood,
		domain.ItemElixir,
		domain.ItemScroll,
		domain.ItemWeapon,
		domain.ItemArmor,
		domain.ItemTreasure,
	}

	itemIndex := 0
	for _, kind := range itemKinds {
		if kind == domain.ItemTreasure {
			totalTreasure := player.Backpack.GetTotalTreasure()
			if totalTreasure > 0 {
				if line >= rect.H-2 {
					return
				}
				w.drawText(rect, line, "Treasure: %d", totalTreasure)
				line++
			}
			continue
		}

		stacks := player.Backpack.GetStacks(kind)
		if len(stacks) == 0 {
			continue
		}

		for _, stack := range stacks {
			if line >= rect.H-2 {
				return
			}

			item := stack.Peek()
			if item == nil {
				continue
			}

			itemName := item.Name()
			properties := w.formatItemProperties(item)
			count := stack.Count

			isSelected := false
			if selection.Model != nil && selection.Model.IsActive() {
				items := selection.Model.Items()
				if itemIndex < len(items) {
					if itemInfo, ok := items[itemIndex].(*service.ItemSelectionItem); ok {
						if itemInfo.Item == item {
							isSelected = selection.Model.SelectedIndex() == itemIndex
						}
					}
				}
			}

			var text string
			if count > 1 {
				text = fmt.Sprintf("%s %s (x%d)", itemName, properties, count)
			} else {
				text = fmt.Sprintf("%s %s", itemName, properties)
			}

			if isSelected {
				w.drawTextHighlighted(rect, line, text)
			} else {
				w.drawText(rect, line, text)
			}
			line++

			if selection.Model != nil && selection.Model.IsActive() {
				items := selection.Model.Items()
				if itemIndex < len(items) {
					if itemInfo, ok := items[itemIndex].(*service.ItemSelectionItem); ok && itemInfo.Item == item {
						itemIndex++
					}
				}
			} else {
				itemIndex++
			}
		}
	}
}

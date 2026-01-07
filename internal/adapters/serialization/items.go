package serialization

import (
	"encoding/json"
	"fmt"
	"sync"

	"rogue/internal/domain"
)

type itemUnmarshaler func([]byte) (domain.Item, error)

var (
	itemRegistry   = make(map[string]itemUnmarshaler)
	itemTypeMap    = make(map[domain.ItemKind]string)
	itemRegistryMu sync.RWMutex
)

func init() {
	registerItem("food", domain.ItemFood, func(data []byte) (domain.Item, error) {
		var food domain.Food
		if err := json.Unmarshal(data, &food); err != nil {
			return nil, err
		}
		return &food, nil
	})
	registerItem("elixir", domain.ItemElixir, func(data []byte) (domain.Item, error) {
		var elixir domain.Elixir
		if err := json.Unmarshal(data, &elixir); err != nil {
			return nil, err
		}
		return &elixir, nil
	})
	registerItem("scroll", domain.ItemScroll, func(data []byte) (domain.Item, error) {
		var scroll domain.Scroll
		if err := json.Unmarshal(data, &scroll); err != nil {
			return nil, err
		}
		return &scroll, nil
	})
	registerItem("weapon", domain.ItemWeapon, func(data []byte) (domain.Item, error) {
		var weapon domain.Weapon
		if err := json.Unmarshal(data, &weapon); err != nil {
			return nil, err
		}
		return &weapon, nil
	})
	registerItem("armor", domain.ItemArmor, func(data []byte) (domain.Item, error) {
		var armor domain.Armor
		if err := json.Unmarshal(data, &armor); err != nil {
			return nil, err
		}
		return &armor, nil
	})
	registerItem("treasure", domain.ItemTreasure, func(data []byte) (domain.Item, error) {
		var treasure domain.Treasure
		if err := json.Unmarshal(data, &treasure); err != nil {
			return nil, err
		}
		return &treasure, nil
	})
}

func registerItem(typeName string, itemKind domain.ItemKind, unmarshaler itemUnmarshaler) {
	itemRegistryMu.Lock()
	defer itemRegistryMu.Unlock()
	itemRegistry[typeName] = unmarshaler
	itemTypeMap[itemKind] = typeName
}

func (s *Serializer) marshalItem(item domain.Item) serializableItem {
	if item == nil {
		return serializableItem{Type: "null", Data: nil}
	}

	itemRegistryMu.RLock()
	typeName, ok := itemTypeMap[item.Type()]
	itemRegistryMu.RUnlock()

	if !ok {
		return serializableItem{Type: "unknown", Data: nil}
	}

	data, err := json.Marshal(item)
	if err != nil {
		return serializableItem{Type: "error", Data: nil}
	}

	return serializableItem{
		Type: typeName,
		Data: data,
	}
}

func (s *Serializer) unmarshalItem(si serializableItem) (domain.Item, error) {
	if si.Type == "null" || si.Type == "error" || si.Type == "unknown" {
		return nil, nil
	}

	itemRegistryMu.RLock()
	unmarshaler, ok := itemRegistry[si.Type]
	itemRegistryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown item type: %s", si.Type)
	}

	return unmarshaler(si.Data)
}

func (s *Serializer) marshalBackpack(backpack *domain.Backpack) *serializableBackpack {
	if backpack == nil {
		return nil
	}

	stacks := make(map[string][]*serializableItemStack)
	allStacks := backpack.GetAllStacks()
	for kind, itemStacks := range allStacks {
		kindStr := itemKindToString(kind)
		serializableStacks := make([]*serializableItemStack, 0, len(itemStacks))
		for _, stack := range itemStacks {
			serializableStacks = append(serializableStacks, &serializableItemStack{
				Item:  s.marshalItem(stack.Item),
				Count: stack.Count,
			})
		}
		stacks[kindStr] = serializableStacks
	}

	return &serializableBackpack{Stacks: stacks}
}

func (s *Serializer) unmarshalBackpack(sb *serializableBackpack) (*domain.Backpack, error) {
	if sb == nil {
		return domain.NewBackpack(), nil
	}

	backpack := domain.NewBackpack()
	for kindStr, serializableStacks := range sb.Stacks {
		kind, err := stringToItemKind(kindStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse item kind: %w", err)
		}

		for _, sStack := range serializableStacks {
			item, err := s.unmarshalItem(sStack.Item)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal item: %w", err)
			}
			if item == nil {
				continue
			}

			stack := domain.NewItemStack(item)
			for i := 1; i < sStack.Count; i++ {
				stack.Add(item)
			}
			backpack.AddStack(kind, stack)
		}
	}

	return backpack, nil
}

func (s *Serializer) marshalEquipment(equipment *domain.Equipment) *serializableEquipment {
	if equipment == nil {
		return nil
	}

	parts := make(map[string]serializableItem)
	allParts := equipment.GetAllParts()
	for part, item := range allParts {
		if item != nil {
			parts[actorPartToString(part)] = s.marshalItem(item)
		}
	}

	return &serializableEquipment{Parts: parts}
}

func (s *Serializer) unmarshalEquipment(se *serializableEquipment) (*domain.Equipment, error) {
	if se == nil {
		return domain.NewEquipment(), nil
	}

	equipment := domain.NewEquipment()
	for partStr, si := range se.Parts {
		part, err := stringToActorPart(partStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse actor part: %w", err)
		}

		item, err := s.unmarshalItem(si)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item for part %s: %w", partStr, err)
		}
		if item != nil {
			equipment.Set(part, item)
		}
	}

	return equipment, nil
}

func itemKindToString(kind domain.ItemKind) string {
	switch kind {
	case domain.ItemFood:
		return "food"
	case domain.ItemElixir:
		return "elixir"
	case domain.ItemScroll:
		return "scroll"
	case domain.ItemWeapon:
		return "weapon"
	case domain.ItemArmor:
		return "armor"
	case domain.ItemTreasure:
		return "treasure"
	default:
		return "unknown"
	}
}

func stringToItemKind(s string) (domain.ItemKind, error) {
	switch s {
	case "food":
		return domain.ItemFood, nil
	case "elixir":
		return domain.ItemElixir, nil
	case "scroll":
		return domain.ItemScroll, nil
	case "weapon":
		return domain.ItemWeapon, nil
	case "armor":
		return domain.ItemArmor, nil
	case "treasure":
		return domain.ItemTreasure, nil
	default:
		return 0, fmt.Errorf("unknown item kind: %s", s)
	}
}

func actorPartToString(part domain.ActorPart) string {
	switch part {
	case domain.ActorPartHead:
		return "head"
	case domain.ActorPartBody:
		return "body"
	case domain.ActorPartHand:
		return "hand"
	case domain.ActorPartLegs:
		return "legs"
	default:
		return "unknown"
	}
}

func stringToActorPart(s string) (domain.ActorPart, error) {
	switch s {
	case "head":
		return domain.ActorPartHead, nil
	case "body":
		return domain.ActorPartBody, nil
	case "hand":
		return domain.ActorPartHand, nil
	case "legs":
		return domain.ActorPartLegs, nil
	default:
		return 0, fmt.Errorf("unknown actor part: %s", s)
	}
}

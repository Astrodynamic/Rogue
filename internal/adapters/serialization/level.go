package serialization

import (
	"fmt"

	"rogue/internal/domain"
)

func (s *Serializer) marshalLevel(level *domain.Level) *serializableLevel {
	if level == nil {
		return nil
	}

	sl := &serializableLevel{
		Rect:      level.Rect,
		Tiles:     level.Tiles,
		Rooms:     level.Rooms,
		Corridors: level.Corridors,
		Exits:     level.Exits,
		Items:     make([]serializableLevelItem, 0),
		Enemies:   make([]serializableLevelEnemy, 0),
	}

	for pos, item := range level.Items {
		si := s.marshalItem(item)
		sl.Items = append(sl.Items, serializableLevelItem{
			Point: pos,
			Item:  si,
		})
	}

	for pos, enemy := range level.Enemies {
		se := s.marshalEnemy(enemy)
		sl.Enemies = append(sl.Enemies, serializableLevelEnemy{
			Point: pos,
			Enemy: se,
		})
	}

	return sl
}

func (s *Serializer) unmarshalLevel(sl *serializableLevel) (*domain.Level, error) {
	if sl == nil {
		return nil, nil
	}

	level := &domain.Level{
		Rect:      sl.Rect,
		Tiles:     sl.Tiles,
		Rooms:     sl.Rooms,
		Corridors: sl.Corridors,
		Exits:     sl.Exits,
		Items:     make(map[domain.Point]domain.Item),
		Enemies:   make(map[domain.Point]domain.Enemy),
	}

	for _, slItem := range sl.Items {
		item, err := s.unmarshalItem(slItem.Item)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item at %v: %w", slItem.Point, err)
		}
		level.Items[slItem.Point] = item
	}

	for _, slEnemy := range sl.Enemies {
		enemy, err := s.unmarshalEnemy(slEnemy.Enemy)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal enemy at %v: %w", slEnemy.Point, err)
		}
		if enemy.GetActor() != nil {
			enemy.GetActor().Point = slEnemy.Point
		}
		level.Enemies[slEnemy.Point] = enemy
	}

	return level, nil
}

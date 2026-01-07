package serialization

import (
	"encoding/json"
	"fmt"

	"rogue/internal/domain"
)

type Serializer struct{}

func NewSerializer() *Serializer {
	return &Serializer{}
}

func (s *Serializer) MarshalWorld(world *domain.World) ([]byte, error) {
	sw := &serializableWorld{
		Level:     s.marshalLevel(world.Level),
		Player:    s.marshalPlayer(world.Player),
		GameState: world.GameState,
	}
	return json.MarshalIndent(sw, "", "  ")
}

func (s *Serializer) UnmarshalWorld(data []byte) (*domain.World, error) {
	var sw serializableWorld
	if err := json.Unmarshal(data, &sw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal world: %w", err)
	}

	level, err := s.unmarshalLevel(sw.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal level: %w", err)
	}

	player, err := s.unmarshalPlayer(sw.Player)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal player: %w", err)
	}

	return &domain.World{
		Level:     level,
		Player:    player,
		GameState: sw.GameState,
	}, nil
}

type serializableWorld struct {
	Level     *serializableLevel  `json:"level"`
	Player    *serializablePlayer `json:"player"`
	GameState *domain.GameState   `json:"game_state"`
}

type serializableLevel struct {
	Rect      domain.Rect              `json:"rect"`
	Tiles     [][]domain.Tile          `json:"tiles"`
	Rooms     []domain.Room            `json:"rooms"`
	Corridors []domain.Corridor        `json:"corridors"`
	Exits     []domain.Point           `json:"exits"`
	Items     []serializableLevelItem  `json:"items"`
	Enemies   []serializableLevelEnemy `json:"enemies"`
}

type serializableLevelItem struct {
	Point domain.Point     `json:"point"`
	Item  serializableItem `json:"item"`
}

type serializableLevelEnemy struct {
	Point domain.Point      `json:"point"`
	Enemy serializableEnemy `json:"enemy"`
}

type serializablePlayer struct {
	Actor     *serializableActor     `json:"actor"`
	Equipment *serializableEquipment `json:"equipment"`
}

type serializableActor struct {
	Point    domain.Point          `json:"point"`
	Stats    domain.Stats          `json:"stats"`
	Name     string                `json:"name"`
	State    domain.ActorState     `json:"state"`
	Effects  []serializableEffect  `json:"effects"`
	Backpack *serializableBackpack `json:"backpack"`
}

type serializableBackpack struct {
	Stacks map[string][]*serializableItemStack `json:"stacks"`
}

type serializableItemStack struct {
	Item  serializableItem `json:"item"`
	Count int              `json:"count"`
}

type serializableEquipment struct {
	Parts map[string]serializableItem `json:"parts"`
}

type serializableItem struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type serializableEnemy struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type serializableEffect struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

package serialization

import (
	"fmt"

	"rogue/internal/domain"
)

func (s *Serializer) marshalPlayer(player *domain.Player) *serializablePlayer {
	if player == nil {
		return nil
	}

	return &serializablePlayer{
		Actor:     s.marshalActor(&player.Actor),
		Equipment: s.marshalEquipment(player.Equipment),
	}
}

func (s *Serializer) unmarshalPlayer(sp *serializablePlayer) (*domain.Player, error) {
	if sp == nil {
		return nil, nil
	}

	actor, err := s.unmarshalActor(sp.Actor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal actor: %w", err)
	}

	equipment, err := s.unmarshalEquipment(sp.Equipment)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal equipment: %w", err)
	}

	return &domain.Player{
		Actor:     *actor,
		Equipment: equipment,
	}, nil
}

func (s *Serializer) marshalActor(actor *domain.Actor) *serializableActor {
	if actor == nil {
		return nil
	}

	effects := make([]serializableEffect, 0, len(actor.Effects))
	for _, effect := range actor.Effects {
		effects = append(effects, s.marshalEffect(effect))
	}

	return &serializableActor{
		Point:    actor.Point,
		Stats:    actor.Stats,
		Name:     actor.Name,
		State:    actor.State,
		Effects:  effects,
		Backpack: s.marshalBackpack(actor.Backpack),
	}
}

func (s *Serializer) unmarshalActor(sa *serializableActor) (*domain.Actor, error) {
	if sa == nil {
		return nil, nil
	}

	effects := make([]domain.Effect, 0, len(sa.Effects))
	for _, se := range sa.Effects {
		effect, err := s.unmarshalEffect(se)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal effect: %w", err)
		}
		effects = append(effects, effect)
	}

	backpack, err := s.unmarshalBackpack(sa.Backpack)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal backpack: %w", err)
	}

	name := sa.Name
	if name == "" {
		name = "Player"
	}
	return &domain.Actor{
		Point:    sa.Point,
		Stats:    sa.Stats,
		Name:     name,
		State:    sa.State,
		Effects:  effects,
		Backpack: backpack,
	}, nil
}

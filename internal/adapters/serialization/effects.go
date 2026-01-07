package serialization

import (
	"encoding/json"
	"fmt"
	"sync"

	"rogue/internal/domain"
)

type effectUnmarshaler func([]byte) (domain.Effect, error)

var (
	effectRegistry   = make(map[string]effectUnmarshaler)
	effectRegistryMu sync.RWMutex
)

func init() {
	registerEffect("health", func(data []byte) (domain.Effect, error) {
		var effect domain.StatEffect
		if err := json.Unmarshal(data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	})
	registerEffect("max_health", func(data []byte) (domain.Effect, error) {
		var effect domain.StatEffect
		if err := json.Unmarshal(data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	})
	registerEffect("dexterity", func(data []byte) (domain.Effect, error) {
		var effect domain.StatEffect
		if err := json.Unmarshal(data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	})
	registerEffect("strength", func(data []byte) (domain.Effect, error) {
		var effect domain.StatEffect
		if err := json.Unmarshal(data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	})
	registerEffect("stat", func(data []byte) (domain.Effect, error) {
		var effect domain.StatEffect
		if err := json.Unmarshal(data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	})
	registerEffect("regeneration", func(data []byte) (domain.Effect, error) {
		var effect domain.RegenEffect
		if err := json.Unmarshal(data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	})
}

func registerEffect(typeName string, unmarshaler effectUnmarshaler) {
	effectRegistryMu.Lock()
	defer effectRegistryMu.Unlock()
	effectRegistry[typeName] = unmarshaler
}

func (s *Serializer) marshalEffect(effect domain.Effect) serializableEffect {
	if effect == nil {
		return serializableEffect{Type: "null", Data: nil}
	}

	var data []byte
	var err error
	var effectType string

	switch v := effect.(type) {
	case *domain.StatEffect:
		if v.Stats.Health != 0 && v.Stats.MaxHealth == 0 && v.Stats.Dexterity == 0 && v.Stats.Strength == 0 {
			effectType = "health"
		} else if v.Stats.MaxHealth != 0 {
			effectType = "max_health"
		} else if v.Stats.Dexterity != 0 {
			effectType = "dexterity"
		} else if v.Stats.Strength != 0 {
			effectType = "strength"
		} else {
			effectType = "stat"
		}
		data, err = json.Marshal(v)
	case *domain.RegenEffect:
		effectType = "regeneration"
		data, err = json.Marshal(v)
	default:
		return serializableEffect{Type: "unknown", Data: nil}
	}

	if err != nil {
		return serializableEffect{Type: "error", Data: nil}
	}

	return serializableEffect{
		Type: effectType,
		Data: data,
	}
}

func (s *Serializer) unmarshalEffect(se serializableEffect) (domain.Effect, error) {
	if se.Type == "null" || se.Type == "error" || se.Type == "unknown" || se.Type == "sleep" {
		return nil, nil
	}

	effectRegistryMu.RLock()
	unmarshaler, ok := effectRegistry[se.Type]
	effectRegistryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown effect type: %s", se.Type)
	}

	return unmarshaler(se.Data)
}

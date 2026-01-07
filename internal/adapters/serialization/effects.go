package serialization

import (
	"encoding/json"
	"fmt"

	"rogue/internal/domain"
)

func (s *Serializer) marshalEffect(effect domain.Effect) serializableEffect {
	if effect == nil {
		return serializableEffect{Type: "null", Data: nil}
	}

	var data []byte
	var err error
	var effectType string

	switch v := effect.(type) {
	case *domain.HealthEffect:
		effectType = "health"
		data, err = json.Marshal(v)
	case *domain.MaxHealthEffect:
		effectType = "max_health"
		data, err = json.Marshal(v)
	case *domain.DexterityEffect:
		effectType = "dexterity"
		data, err = json.Marshal(v)
	case *domain.StrengthEffect:
		effectType = "strength"
		data, err = json.Marshal(v)
	case *domain.SleepEffect:
		effectType = "sleep"
		data, err = json.Marshal(v)
	case *domain.RegenerationEffect:
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
	switch se.Type {
	case "null", "error", "unknown":
		return nil, nil
	case "health":
		var effect domain.HealthEffect
		if err := json.Unmarshal(se.Data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	case "max_health":
		var effect domain.MaxHealthEffect
		if err := json.Unmarshal(se.Data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	case "dexterity":
		var effect domain.DexterityEffect
		if err := json.Unmarshal(se.Data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	case "strength":
		var effect domain.StrengthEffect
		if err := json.Unmarshal(se.Data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	case "sleep":
		var effect domain.SleepEffect
		if err := json.Unmarshal(se.Data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	case "regeneration":
		var effect domain.RegenerationEffect
		if err := json.Unmarshal(se.Data, &effect); err != nil {
			return nil, err
		}
		return &effect, nil
	default:
		return nil, fmt.Errorf("unknown effect type: %s", se.Type)
	}
}

package serialization

import (
	"encoding/json"
	"fmt"

	"rogue/internal/domain"
)

func (s *Serializer) marshalEnemy(enemy domain.Enemy) serializableEnemy {
	if enemy == nil {
		return serializableEnemy{Type: "null", Data: nil}
	}

	var data []byte
	var err error
	var enemyType string

	switch v := enemy.(type) {
	case *domain.Zombie:
		enemyType = "zombie"
		data, err = json.Marshal(v)
	case *domain.Vampire:
		enemyType = "vampire"
		data, err = json.Marshal(v)
	case *domain.Ghost:
		enemyType = "ghost"
		data, err = json.Marshal(v)
	case *domain.Ogre:
		enemyType = "ogre"
		data, err = json.Marshal(v)
	case *domain.SnakeMage:
		enemyType = "snake_mage"
		data, err = json.Marshal(v)
	case *domain.Mimic:
		enemyType = "mimic"
		data, err = json.Marshal(v)
	default:
		return serializableEnemy{Type: "unknown", Data: nil}
	}

	if err != nil {
		return serializableEnemy{Type: "error", Data: nil}
	}

	return serializableEnemy{
		Type: enemyType,
		Data: data,
	}
}

func (s *Serializer) unmarshalEnemy(se serializableEnemy) (domain.Enemy, error) {
	switch se.Type {
	case "null", "error", "unknown":
		return nil, nil
	case "zombie":
		var zombie domain.Zombie
		if err := json.Unmarshal(se.Data, &zombie); err != nil {
			return nil, err
		}
		return &zombie, nil
	case "vampire":
		var vampire domain.Vampire
		if err := json.Unmarshal(se.Data, &vampire); err != nil {
			return nil, err
		}
		return &vampire, nil
	case "ghost":
		var ghost domain.Ghost
		if err := json.Unmarshal(se.Data, &ghost); err != nil {
			return nil, err
		}
		return &ghost, nil
	case "ogre":
		var ogre domain.Ogre
		if err := json.Unmarshal(se.Data, &ogre); err != nil {
			return nil, err
		}
		return &ogre, nil
	case "snake_mage":
		var snakeMage domain.SnakeMage
		if err := json.Unmarshal(se.Data, &snakeMage); err != nil {
			return nil, err
		}
		return &snakeMage, nil
	case "mimic":
		var mimic domain.Mimic
		if err := json.Unmarshal(se.Data, &mimic); err != nil {
			return nil, err
		}
		return &mimic, nil
	default:
		return nil, fmt.Errorf("unknown enemy type: %s", se.Type)
	}
}

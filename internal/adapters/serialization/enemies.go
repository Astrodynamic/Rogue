package serialization

import (
	"encoding/json"
	"fmt"
	"sync"

	"rogue/internal/domain"
)

type enemyMarshaler func(domain.Enemy) ([]byte, error)
type enemyUnmarshaler func([]byte) (domain.Enemy, error)

var (
	enemyRegistry     = make(map[string]enemyUnmarshaler)
	enemyTypeRegistry = make(map[domain.EnemyType]string)
	enemyRegistryMu   sync.RWMutex
)

func init() {
	registerEnemy("zombie", domain.EnemyTypeZombie, func(data []byte) (domain.Enemy, error) {
		var zombie domain.Zombie
		if err := json.Unmarshal(data, &zombie); err != nil {
			return nil, err
		}
		return &zombie, nil
	})
	registerEnemy("vampire", domain.EnemyTypeVampire, func(data []byte) (domain.Enemy, error) {
		var vampire domain.Vampire
		if err := json.Unmarshal(data, &vampire); err != nil {
			return nil, err
		}
		return &vampire, nil
	})
	registerEnemy("ghost", domain.EnemyTypeGhost, func(data []byte) (domain.Enemy, error) {
		var ghost domain.Ghost
		if err := json.Unmarshal(data, &ghost); err != nil {
			return nil, err
		}
		return &ghost, nil
	})
	registerEnemy("ogre", domain.EnemyTypeOgre, func(data []byte) (domain.Enemy, error) {
		var ogre domain.Ogre
		if err := json.Unmarshal(data, &ogre); err != nil {
			return nil, err
		}
		return &ogre, nil
	})
	registerEnemy("snake_mage", domain.EnemyTypeSnakeMage, func(data []byte) (domain.Enemy, error) {
		var snakeMage domain.SnakeMage
		if err := json.Unmarshal(data, &snakeMage); err != nil {
			return nil, err
		}
		return &snakeMage, nil
	})
	registerEnemy("mimic", domain.EnemyTypeMimic, func(data []byte) (domain.Enemy, error) {
		var mimic domain.Mimic
		if err := json.Unmarshal(data, &mimic); err != nil {
			return nil, err
		}
		return &mimic, nil
	})
}

func registerEnemy(typeName string, enemyType domain.EnemyType, unmarshaler enemyUnmarshaler) {
	enemyRegistryMu.Lock()
	defer enemyRegistryMu.Unlock()
	enemyRegistry[typeName] = unmarshaler
	enemyTypeRegistry[enemyType] = typeName
}

func (s *Serializer) marshalEnemy(enemy domain.Enemy) serializableEnemy {
	if enemy == nil {
		return serializableEnemy{Type: "null", Data: nil}
	}

	enemyRegistryMu.RLock()
	typeName, ok := enemyTypeRegistry[enemy.GetEnemyType()]
	enemyRegistryMu.RUnlock()

	if !ok {
		return serializableEnemy{Type: "unknown", Data: nil}
	}

	data, err := json.Marshal(enemy)
	if err != nil {
		return serializableEnemy{Type: "error", Data: nil}
	}

	return serializableEnemy{
		Type: typeName,
		Data: data,
	}
}

func (s *Serializer) unmarshalEnemy(se serializableEnemy) (domain.Enemy, error) {
	if se.Type == "null" || se.Type == "error" || se.Type == "unknown" {
		return nil, nil
	}

	enemyRegistryMu.RLock()
	unmarshaler, ok := enemyRegistry[se.Type]
	enemyRegistryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown enemy type: %s", se.Type)
	}

	return unmarshaler(se.Data)
}

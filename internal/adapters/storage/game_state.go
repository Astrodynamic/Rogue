package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"rogue/internal/adapters/serialization"
	"rogue/internal/domain"
)

const gameStateFileName = "game_state.json"

type GameStateStorage struct {
	filePath   string
	serializer *serialization.Serializer
}

func NewGameStateStorage(savesDir string) *GameStateStorage {
	return &GameStateStorage{
		filePath:   filepath.Join(savesDir, gameStateFileName),
		serializer: serialization.NewSerializer(),
	}
}

func (s *GameStateStorage) Save(world *domain.World) error {
	data, err := s.serializer.MarshalWorld(world)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create saves directory: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write game state file: %w", err)
	}

	return nil
}

func (s *GameStateStorage) Load() (*domain.World, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read game state file: %w", err)
	}

	world, err := s.serializer.UnmarshalWorld(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal game state: %w", err)
	}

	return world, nil
}

func (s *GameStateStorage) HasSave() bool {
	_, err := os.Stat(s.filePath)
	return err == nil
}

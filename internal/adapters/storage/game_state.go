package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"rogue/internal/adapters/serialization"
	"rogue/internal/domain"
)

type StateStore struct {
	savesDir   string
	serializer *serialization.Serializer
}

func NewStateStore(savesDir string) *StateStore {
	return &StateStore{
		savesDir:   savesDir,
		serializer: serialization.NewSerializer(),
	}
}

func (s *StateStore) getFilePath(playerName string) string {
	sanitized := sanitizePlayerName(playerName)
	fileName := fmt.Sprintf("%s_game_state.json", sanitized)
	return filepath.Join(s.savesDir, fileName)
}

func (s *StateStore) Save(world *domain.World) error {
	if world == nil || world.GameState == nil {
		return fmt.Errorf("world or game state is nil")
	}

	data, err := s.serializer.MarshalWorld(world)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	filePath := s.getFilePath(world.Player.Name)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create saves directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write game state file: %w", err)
	}

	return nil
}

func (s *StateStore) Load(playerName string) (*domain.World, error) {
	if playerName == "" {
		return nil, fmt.Errorf("player name cannot be empty")
	}

	filePath := s.getFilePath(playerName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read game state file: %w", err)
	}

	world, err := s.serializer.UnmarshalWorld(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal world: %w", err)
	}

	return world, nil
}

func (s *StateStore) HasSave(playerName string) bool {
	if playerName == "" {
		return false
	}
	filePath := s.getFilePath(playerName)
	_, err := os.Stat(filePath)
	return err == nil
}

func (s *StateStore) ListSaves() ([]string, error) {
	if err := os.MkdirAll(s.savesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create saves directory: %w", err)
	}

	entries, err := os.ReadDir(s.savesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read saves directory: %w", err)
	}

	var playerNames []string
	seen := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, "_game_state.json") {
			playerName := strings.TrimSuffix(name, "_game_state.json")
			if !seen[playerName] {
				playerNames = append(playerNames, playerName)
				seen[playerName] = true
			}
		}
	}

	return playerNames, nil
}

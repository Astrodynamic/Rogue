package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"rogue/internal/domain"
)

type StatsStore struct {
	savesDir string
}

func NewStatsStore(savesDir string) *StatsStore {
	return &StatsStore{
		savesDir: savesDir,
	}
}

func (s *StatsStore) getFilePath(playerName string) string {
	sanitized := sanitizePlayerName(playerName)
	fileName := fmt.Sprintf("%s_statistics.json", sanitized)
	return filepath.Join(s.savesDir, fileName)
}

func (s *StatsStore) SavePlaythrough(stats *domain.PlayStats) error {
	if stats == nil {
		return fmt.Errorf("statistics cannot be nil")
	}

	playerName := stats.PlayerName
	if playerName == "" {
		playerName = "Player"
	}

	playthroughs, err := s.loadPlayerPlaythroughs(playerName)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load existing playthroughs: %w", err)
	}

	playthroughs = append(playthroughs, stats)

	data, err := json.MarshalIndent(playthroughs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal statistics: %w", err)
	}

	filePath := s.getFilePath(playerName)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create saves directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write statistics file: %w", err)
	}

	return nil
}

func (s *StatsStore) loadPlayerPlaythroughs(playerName string) ([]*domain.PlayStats, error) {
	filePath := s.getFilePath(playerName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var playthroughs []*domain.PlayStats
	if err := json.Unmarshal(data, &playthroughs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal statistics: %w", err)
	}

	return playthroughs, nil
}

func (s *StatsStore) LoadAllPlaythroughs() ([]*domain.PlayStats, error) {
	if err := os.MkdirAll(s.savesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create saves directory: %w", err)
	}

	entries, err := os.ReadDir(s.savesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read saves directory: %w", err)
	}

	var allPlaythroughs []*domain.PlayStats
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, "_statistics.json") {
			playerName := strings.TrimSuffix(name, "_statistics.json")
			playthroughs, err := s.loadPlayerPlaythroughs(playerName)
			if err != nil {
				if !os.IsNotExist(err) {
					return nil, fmt.Errorf("failed to load statistics for player %s: %w", playerName, err)
				}
				continue
			}
			allPlaythroughs = append(allPlaythroughs, playthroughs...)
		}
	}

	return allPlaythroughs, nil
}

func (s *StatsStore) GetLeaderboard(limit int) ([]*domain.PlayStats, error) {
	playthroughs, err := s.LoadAllPlaythroughs()
	if err != nil {
		if os.IsNotExist(err) {
			return []*domain.PlayStats{}, nil
		}
		return nil, err
	}

	sort.Slice(playthroughs, func(i, j int) bool {
		if playthroughs[i].TreasureCollected != playthroughs[j].TreasureCollected {
			return playthroughs[i].TreasureCollected > playthroughs[j].TreasureCollected
		}
		return playthroughs[i].DeepestLevel > playthroughs[j].DeepestLevel
	})

	if limit > 0 && limit < len(playthroughs) {
		return playthroughs[:limit], nil
	}

	return playthroughs, nil
}

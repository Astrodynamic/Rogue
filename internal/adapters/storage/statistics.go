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

type StatisticsStorage struct {
	savesDir string
}

func NewStatisticsStorage(savesDir string) *StatisticsStorage {
	return &StatisticsStorage{
		savesDir: savesDir,
	}
}

func (s *StatisticsStorage) getFilePath(playerName string) string {
	sanitized := sanitizePlayerName(playerName)
	fileName := fmt.Sprintf("%s_statistics.json", sanitized)
	return filepath.Join(s.savesDir, fileName)
}

func (s *StatisticsStorage) SavePlaythrough(stats *domain.PlaythroughStatistics) error {
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

func (s *StatisticsStorage) loadPlayerPlaythroughs(playerName string) ([]*domain.PlaythroughStatistics, error) {
	filePath := s.getFilePath(playerName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var playthroughs []*domain.PlaythroughStatistics
	if err := json.Unmarshal(data, &playthroughs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal statistics: %w", err)
	}

	return playthroughs, nil
}

func (s *StatisticsStorage) LoadAllPlaythroughs() ([]*domain.PlaythroughStatistics, error) {
	if err := os.MkdirAll(s.savesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create saves directory: %w", err)
	}

	entries, err := os.ReadDir(s.savesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read saves directory: %w", err)
	}

	var allPlaythroughs []*domain.PlaythroughStatistics
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

func (s *StatisticsStorage) GetLeaderboard(limit int) ([]*domain.PlaythroughStatistics, error) {
	playthroughs, err := s.LoadAllPlaythroughs()
	if err != nil {
		if os.IsNotExist(err) {
			return []*domain.PlaythroughStatistics{}, nil
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

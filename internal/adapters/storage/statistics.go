package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"rogue/internal/domain"
)

const statisticsFileName = "statistics.json"

type StatisticsStorage struct {
	filePath string
}

func NewStatisticsStorage(savesDir string) *StatisticsStorage {
	return &StatisticsStorage{
		filePath: filepath.Join(savesDir, statisticsFileName),
	}
}

func (s *StatisticsStorage) SavePlaythrough(stats *domain.PlaythroughStatistics) error {
	playthroughs, err := s.LoadAllPlaythroughs()
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to load existing playthroughs: %w", err)
	}

	playthroughs = append(playthroughs, stats)

	data, err := json.MarshalIndent(playthroughs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal statistics: %w", err)
	}

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create saves directory: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write statistics file: %w", err)
	}

	return nil
}

func (s *StatisticsStorage) LoadAllPlaythroughs() ([]*domain.PlaythroughStatistics, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var playthroughs []*domain.PlaythroughStatistics
	if err := json.Unmarshal(data, &playthroughs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal statistics: %w", err)
	}

	return playthroughs, nil
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

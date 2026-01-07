package service

import "rogue/internal/domain"

type GameStateStorage interface {
	Save(world *domain.World) error
	Load(playerName string) (*domain.World, error)
	HasSave(playerName string) bool
	ListSaves() ([]string, error)
}

type StatisticsStorage interface {
	SavePlaythrough(stats *domain.PlaythroughStatistics) error
	GetLeaderboard(limit int) ([]*domain.PlaythroughStatistics, error)
}

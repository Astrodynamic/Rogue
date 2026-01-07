package service

import "rogue/internal/domain"

type StateStore interface {
	Save(world *domain.World) error
	Load(playerName string) (*domain.World, error)
	HasSave(playerName string) bool
	ListSaves() ([]string, error)
}

type StatsStore interface {
	SavePlaythrough(stats *domain.PlayStats) error
	GetLeaderboard(limit int) ([]*domain.PlayStats, error)
}

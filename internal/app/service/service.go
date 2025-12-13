package service

import "github.com/Astrodynamic/Rogue/internal/domain/ports"

type Service struct {
	leaderboard ports.LeaderboardRepository
	sessions    ports.SessionRepository
}

func New(leaderboard ports.LeaderboardRepository, sessions ports.SessionRepository) *Service {
	return &Service{leaderboard: leaderboard, sessions: sessions}
}

func (s *Service) RecordRun(r ports.RunResult) error { return s.leaderboard.Record(r) }

func (s *Service) TopRuns(n int) ([]ports.RunResult, error) { return s.leaderboard.Top(n) }

func (s *Service) SaveSession(slot int, sess ports.SavedSession) error {
	return s.sessions.SaveSlot(slot, sess)
}

func (s *Service) LoadSession(slot int) (ports.SavedSession, bool, error) {
	return s.sessions.LoadSlot(slot)
}

func (s *Service) ClearSession(slot int) error { return s.sessions.ClearSlot(slot) }

func (s *Service) ListSaveSlots() ([]ports.SaveSlot, error) { return s.sessions.ListSlots() }

func (s *Service) AllRuns() ([]ports.RunResult, error) { return s.leaderboard.All() }

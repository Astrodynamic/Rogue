package ports

import "rogue/internal/core/domain/game"

type RunResult struct {
	Name         string
	DeepestLevel int
	Treasure     int
	Stats        game.Stats
}

type LeaderboardRepository interface {
	Record(RunResult) error
	Top(n int) ([]RunResult, error)
	All() ([]RunResult, error)
}

type SavedSession struct {
	Snapshot game.Snapshot
	Name     string
}

type SaveSlot struct {
	Slot      int
	Name      string
	Level     int
	Treasure  int
	UpdatedAt int64
	Empty     bool
}

type SessionRepository interface {
	SaveSlot(slot int, sess SavedSession) error
	LoadSlot(slot int) (SavedSession, bool, error)
	ClearSlot(slot int) error
	ListSlots() ([]SaveSlot, error)
}

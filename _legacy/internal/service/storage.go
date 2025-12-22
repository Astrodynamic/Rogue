package service

import "rogue/internal/domain"

// Storage defines what the service needs from persistence layer.
// In Go, interfaces are defined where they're used (consumer-side).
// Storage implementations (JSON, memory, etc.) live in internal/adapters/storage.
type Storage interface {
	// Leaderboard operations
	RecordRun(r RunResult) error
	AllRuns() ([]RunResult, error)

	// Session operations
	SaveSession(slot int, sess SavedSession) error
	LoadSession(slot int) (SavedSession, bool, error)
	ClearSession(slot int) error
	ListSaveSlots() ([]SaveSlot, error)
}

// RunResult is a completed game run (leaderboard entry).
type RunResult struct {
	Name         string
	DeepestLevel int
	Treasure     int
	Stats        domain.Stats
}

// SavedSession is a saved game state.
type SavedSession struct {
	Snapshot domain.Snapshot
	Name     string
}

// SaveSlot is metadata about a save slot (for UI listing).
type SaveSlot struct {
	Slot      int
	Name      string
	Level     int
	Treasure  int
	UpdatedAt int64
	Empty     bool
}

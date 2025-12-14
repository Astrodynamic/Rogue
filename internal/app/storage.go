package app

import "rogue/internal/game"

// Storage defines what the app needs from persistence layer.
// In Go, interfaces are defined where they're used (consumer-side).
// Storage implementations (JSON, memory, etc.) live in internal/storage.
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
	Stats        game.Stats
}

// SavedSession is a saved game state.
type SavedSession struct {
	Snapshot game.Snapshot
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

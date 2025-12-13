package usecase

import "github.com/Astrodynamic/Rogue/internal/domain/ports"

// App is the usecase-facing contract (application boundary).
// The usecase depends on this interface, not on concrete adapters, to keep clean architecture/DDD:
// - domain is independent
// - usecases orchestrate domain and call ports through App
// - adapters (storage/tui) sit at the edges and are injected from main
type App interface {
	RecordRun(r ports.RunResult) error
	AllRuns() ([]ports.RunResult, error)

	SaveSession(slot int, sess ports.SavedSession) error
	LoadSession(slot int) (ports.SavedSession, bool, error)
	ClearSession(slot int) error
	ListSaveSlots() ([]ports.SaveSlot, error)
}

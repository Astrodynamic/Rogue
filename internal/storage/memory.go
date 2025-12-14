package storage

import (
	"sort"
	"sync"

	"rogue/internal/app"
)

type InMemory struct {
	mu      sync.Mutex
	results []app.RunResult
	slots   map[int]app.SavedSession
}

func NewInMemory() *InMemory { return &InMemory{slots: map[int]app.SavedSession{}} }

func (m *InMemory) RecordRun(r app.RunResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.results = append(m.results, r)
	sort.SliceStable(m.results, func(i, j int) bool {
		return m.results[i].Treasure > m.results[j].Treasure
	})
	return nil
}

func (m *InMemory) TopRuns(n int) ([]app.RunResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if n <= 0 || n > len(m.results) {
		n = len(m.results)
	}
	out := make([]app.RunResult, n)
	copy(out, m.results[:n])
	return out, nil
}

func (m *InMemory) AllRuns() ([]app.RunResult, error) {
	return m.TopRuns(0)
}

func (m *InMemory) SaveSession(slot int, s app.SavedSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if slot < 1 || slot > 3 {
		return nil
	}
	m.slots[slot] = s
	return nil
}

func (m *InMemory) LoadSession(slot int) (app.SavedSession, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.slots[slot]
	return s, ok, nil
}

func (m *InMemory) ClearSession(slot int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.slots, slot)
	return nil
}

func (m *InMemory) ListSaveSlots() ([]app.SaveSlot, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]app.SaveSlot, 0, 3)
	for slot := 1; slot <= 3; slot++ {
		if s, ok := m.slots[slot]; ok {
			out = append(out, app.SaveSlot{Slot: slot, Name: s.Name, Level: s.Snapshot.LevelDepth, Treasure: s.Snapshot.Stats.TotalTreasure, Empty: false})
		} else {
			out = append(out, app.SaveSlot{Slot: slot, Empty: true})
		}
	}
	return out, nil
}

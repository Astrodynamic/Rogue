package storage

import (
	"sort"
	"sync"

	"rogue/internal/service"
)

type InMemory struct {
	mu      sync.Mutex
	results []service.RunResult
	slots   map[int]service.SavedSession
}

func NewInMemory() *InMemory { return &InMemory{slots: map[int]service.SavedSession{}} }

func (m *InMemory) RecordRun(r service.RunResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.results = append(m.results, r)
	sort.SliceStable(m.results, func(i, j int) bool {
		return m.results[i].Treasure > m.results[j].Treasure
	})
	return nil
}

func (m *InMemory) TopRuns(n int) ([]service.RunResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if n <= 0 || n > len(m.results) {
		n = len(m.results)
	}
	out := make([]service.RunResult, n)
	copy(out, m.results[:n])
	return out, nil
}

func (m *InMemory) AllRuns() ([]service.RunResult, error) {
	return m.TopRuns(0)
}

func (m *InMemory) SaveSession(slot int, s service.SavedSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if slot < 1 || slot > 3 {
		return nil
	}
	m.slots[slot] = s
	return nil
}

func (m *InMemory) LoadSession(slot int) (service.SavedSession, bool, error) {
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

func (m *InMemory) ListSaveSlots() ([]service.SaveSlot, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]service.SaveSlot, 0, 3)
	for slot := 1; slot <= 3; slot++ {
		if s, ok := m.slots[slot]; ok {
			out = append(out, service.SaveSlot{Slot: slot, Name: s.Name, Level: s.Snapshot.LevelDepth, Treasure: s.Snapshot.Stats.TotalTreasure, Empty: false})
		} else {
			out = append(out, service.SaveSlot{Slot: slot, Empty: true})
		}
	}
	return out, nil
}

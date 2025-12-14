package storage

import (
	"sort"
	"sync"

	"rogue/internal/domain/ports"
)

type InMemory struct {
	mu      sync.Mutex
	results []ports.RunResult
	slots   map[int]ports.SavedSession
}

func NewInMemory() *InMemory { return &InMemory{slots: map[int]ports.SavedSession{}} }

func (m *InMemory) Record(r ports.RunResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.results = append(m.results, r)
	sort.SliceStable(m.results, func(i, j int) bool {
		return m.results[i].Treasure > m.results[j].Treasure
	})
	return nil
}

func (m *InMemory) Top(n int) ([]ports.RunResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if n <= 0 || n > len(m.results) {
		n = len(m.results)
	}
	out := make([]ports.RunResult, n)
	copy(out, m.results[:n])
	return out, nil
}

func (m *InMemory) All() ([]ports.RunResult, error) {
	return m.Top(0)
}

func (m *InMemory) SaveSlot(slot int, s ports.SavedSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if slot < 1 || slot > 3 {
		return nil
	}
	m.slots[slot] = s
	return nil
}

func (m *InMemory) LoadSlot(slot int) (ports.SavedSession, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.slots[slot]
	return s, ok, nil
}

func (m *InMemory) ClearSlot(slot int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.slots, slot)
	return nil
}

func (m *InMemory) ListSlots() ([]ports.SaveSlot, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]ports.SaveSlot, 0, 3)
	for slot := 1; slot <= 3; slot++ {
		if s, ok := m.slots[slot]; ok {
			out = append(out, ports.SaveSlot{Slot: slot, Name: s.Name, Level: s.Snapshot.LevelDepth, Treasure: s.Snapshot.Stats.TotalTreasure, Empty: false})
		} else {
			out = append(out, ports.SaveSlot{Slot: slot, Empty: true})
		}
	}
	return out, nil
}

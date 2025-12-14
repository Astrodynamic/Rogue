package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"

	"rogue/internal/app"
)

type JSONStore struct {
	mu   sync.Mutex
	path string
	data fileData
}

type fileData struct {
	Saves map[string]app.SavedSession `json:"saves,omitempty"`
	Runs  []app.RunResult             `json:"runs,omitempty"`

	// Backward compatibility: older versions stored a single session under "saved".
	Saved *app.SavedSession `json:"saved,omitempty"`
}

func NewJSON(path string) *JSONStore {
	return &JSONStore{path: path}
}

func (s *JSONStore) loadLocked() error {
	b, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.data = fileData{}
			return nil
		}
		return err
	}
	if len(b) == 0 {
		s.data = fileData{}
		return nil
	}
	var d fileData
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	// Migrate legacy "saved" to slot 1.
	if d.Saves == nil && d.Saved != nil {
		d.Saves = map[string]app.SavedSession{"1": *d.Saved}
		d.Saved = nil
	}
	s.data = d
	return nil
}

func (s *JSONStore) saveLocked() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *JSONStore) RecordRun(r app.RunResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return err
	}
	s.data.Runs = append(s.data.Runs, r)
	sort.SliceStable(s.data.Runs, func(i, j int) bool {
		return s.data.Runs[i].Treasure > s.data.Runs[j].Treasure
	})
	return s.saveLocked()
}

func (s *JSONStore) TopRuns(n int) ([]app.RunResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, err
	}
	if n <= 0 || n > len(s.data.Runs) {
		n = len(s.data.Runs)
	}
	out := make([]app.RunResult, n)
	copy(out, s.data.Runs[:n])
	return out, nil
}

func (s *JSONStore) AllRuns() ([]app.RunResult, error) { return s.TopRuns(0) }

func (s *JSONStore) Save(sess app.SavedSession) error {
	return s.SaveSession(1, sess)
}

func (s *JSONStore) SaveSession(slot int, sess app.SavedSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if slot < 1 || slot > 3 {
		return nil
	}
	if err := s.loadLocked(); err != nil {
		return err
	}
	if s.data.Saves == nil {
		s.data.Saves = map[string]app.SavedSession{}
	}
	key := itoa(slot)
	s.data.Saves[key] = sess
	return s.saveLocked()
}

func (s *JSONStore) LoadSession(slot int) (app.SavedSession, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if slot < 1 || slot > 3 {
		return app.SavedSession{}, false, nil
	}
	if err := s.loadLocked(); err != nil {
		return app.SavedSession{}, false, err
	}
	if s.data.Saves == nil {
		return app.SavedSession{}, false, nil
	}
	key := itoa(slot)
	val, ok := s.data.Saves[key]
	return val, ok, nil
}

func (s *JSONStore) ClearSession(slot int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if slot < 1 || slot > 3 {
		return nil
	}
	if err := s.loadLocked(); err != nil {
		return err
	}
	if s.data.Saves == nil {
		return nil
	}
	delete(s.data.Saves, itoa(slot))
	return s.saveLocked()
}

func (s *JSONStore) ListSaveSlots() ([]app.SaveSlot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.loadLocked(); err != nil {
		return nil, err
	}
	out := make([]app.SaveSlot, 0, 3)
	for slot := 1; slot <= 3; slot++ {
		if s.data.Saves == nil {
			out = append(out, app.SaveSlot{Slot: slot, Empty: true})
			continue
		}
		val, ok := s.data.Saves[itoa(slot)]
		if !ok {
			out = append(out, app.SaveSlot{Slot: slot, Empty: true})
			continue
		}
		out = append(out, app.SaveSlot{
			Slot:     slot,
			Name:     val.Name,
			Level:    val.Snapshot.LevelDepth,
			Treasure: val.Snapshot.Stats.TotalTreasure,
			Empty:    false,
		})
	}
	return out, nil
}

// Backward-compat methods removed from ports; keep tiny helpers local.
func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [16]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + (v % 10))
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

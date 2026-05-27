package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

// Store defines load/save behavior for counter state.
type Store interface {
	Load() (map[string]int, error)
	Save(state map[string]int) error
}

// JSONStore stores counter state in a local JSON file.
type JSONStore struct {
	path   string
	mu     sync.Mutex
	logger *zap.Logger
}

// NewJSONStore creates a JSONStore for the provided file path.
func NewJSONStore(path string, logger *zap.Logger) *JSONStore {
	return &JSONStore{path: path, logger: logger}
}

// Load reads state from the JSON file. If file doesn't exist, it creates it.
func (s *JSONStore) Load() (map[string]int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
		empty := map[string]int{}
		if err := s.saveLocked(empty); err != nil {
			return nil, err
		}
		return empty, nil
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	state := map[string]int{}
	if len(data) == 0 {
		return state, nil
	}

	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return state, nil
}

// Save writes the state into the JSON file safely.
func (s *JSONStore) Save(state map[string]int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.saveLocked(state)
}

func (s *JSONStore) saveLocked(state map[string]int) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	payload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o644); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, s.path); err != nil {
		return err
	}

	if s.logger != nil {
		s.logger.Debug("state saved", zap.String("path", s.path))
	}
	return nil
}

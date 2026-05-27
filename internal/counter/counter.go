package counter

import (
	"errors"
	"sync"

	"distributed-counter-go/internal/storage"

	"go.uber.org/zap"
)

// Counter holds the CRDT G-Counter state for this node.
type Counter struct {
	nodeID string
	mu     sync.Mutex
	state  map[string]int
	store  storage.Store
	logger *zap.Logger
}

// NewCounter creates a new counter and loads state from storage.
func NewCounter(nodeID string, store storage.Store, logger *zap.Logger) (*Counter, error) {
	state, err := store.Load()
	if err != nil {
		return nil, err
	}
	if state == nil {
		state = map[string]int{}
	}
	if _, ok := state[nodeID]; !ok {
		state[nodeID] = 0
		if err := store.Save(state); err != nil {
			return nil, err
		}
	}

	return &Counter{
		nodeID: nodeID,
		state:  state,
		store:  store,
		logger: logger,
	}, nil
}

// Increment increases the local counter and persists the state.
func (c *Counter) Increment(amount int) (int, int, error) {
	if amount <= 0 {
		return 0, 0, errors.New("amount must be greater than zero")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.state[c.nodeID] += amount
	if err := c.store.Save(c.state); err != nil {
		return 0, 0, err
	}

	local := c.state[c.nodeID]
	return local, totalLocked(c.state), nil
}

// MergeState merges incoming CRDT state and persists if changed.
func (c *Counter) MergeState(incoming map[string]int) (int, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	merged, changed := MergeState(c.state, incoming)
	c.state = merged

	if changed {
		if err := c.store.Save(c.state); err != nil {
			return 0, false, err
		}
	}

	return totalLocked(c.state), changed, nil
}

// GetState returns a copy of the current state.
func (c *Counter) GetState() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()

	copyState := make(map[string]int, len(c.state))
	for k, v := range c.state {
		copyState[k] = v
	}

	return copyState
}

// Total returns the current global counter value.
func (c *Counter) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return totalLocked(c.state)
}

func totalLocked(state map[string]int) int {
	total := 0
	for _, v := range state {
		total += v
	}
	return total
}

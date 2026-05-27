package cluster

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"distributed-counter-go/internal/model"

	"go.uber.org/zap"
)

// SyncLoop periodically syncs state with peers for eventual consistency.
func (m *Manager) SyncLoop(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.SyncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.SyncCluster()
		}
	}
}

// SyncCluster pulls state from peers and merges into local counter.
func (m *Manager) SyncCluster() {
	peers := m.GetPeers()
	localState := m.counter.GetState()

	for _, peer := range peers {
		peerState, err := m.syncWithPeer(peer, localState)
		if err != nil {
			if m.logger != nil {
				m.logger.Warn("sync with peer failed", zap.String("peer", peer), zap.Error(err))
			}
			continue
		}
		_, _, _ = m.counter.MergeState(peerState)
	}
}

// AggregateCounter collects state from peers and returns global counter and active nodes.
func (m *Manager) AggregateCounter() (int, int, error) {
	localState := m.counter.GetState()
	merged := localState
	active := 1

	for _, peer := range m.GetPeers() {
		peerState, err := m.syncWithPeer(peer, localState)
		if err != nil {
			continue
		}
		merged, _ = mergeStates(merged, peerState)
		active++
	}

	return total(merged), active, nil
}

// BroadcastIncrement shares the latest local state to peers asynchronously.
func (m *Manager) BroadcastIncrement() {
	go func() {
		localState := m.counter.GetState()
		for _, peer := range m.GetPeers() {
			_, _ = m.syncWithPeer(peer, localState)
		}
	}()
}

func (m *Manager) syncWithPeer(peer string, localState map[string]int) (map[string]int, error) {
	payload := model.SyncRequest{
		State: localState,
		Peer:  m.cfg.NodeURL,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := peer + "/sync"

	var lastErr error
	for i := 0; i < 3; i++ { // initial try + 2 retries
		resp, err := m.client.Post(url, "application/json", bytes.NewReader(data))
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			lastErr = errors.New("sync request failed")
			continue
		}

		var out model.SyncResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			_ = resp.Body.Close()
			lastErr = err
			continue
		}

		_ = resp.Body.Close()
		return out.State, nil
	}

	return nil, lastErr
}

func mergeStates(base, incoming map[string]int) (map[string]int, bool) {
	merged := make(map[string]int, len(base)+len(incoming))
	changed := false
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range incoming {
		if cur, ok := merged[k]; !ok || v > cur {
			merged[k] = v
			changed = true
		}
	}
	return merged, changed
}

func total(state map[string]int) int {
	count := 0
	for _, v := range state {
		count += v
	}
	return count
}

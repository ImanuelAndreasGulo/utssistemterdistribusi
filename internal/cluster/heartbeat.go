package cluster

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Heartbeat checks peer health periodically with retries.
func (m *Manager) Heartbeat(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.checkPeersHealth()
		}
	}
}

func (m *Manager) checkPeersHealth() {
	peers := m.GetPeers()
	for _, peer := range peers {
		healthy := m.checkHealthWithRetry(peer)
		m.SetHealth(peer, healthy)
	}
}

func (m *Manager) checkHealthWithRetry(peer string) bool {
	for i := 0; i < 3; i++ { // initial try + 2 retries
		if m.checkHealth(peer) {
			return true
		}
	}
	return false
}

func (m *Manager) checkHealth(peer string) bool {
	url := peer + "/health"
	resp, err := m.client.Get(url)
	if err != nil {
		if m.logger != nil {
			m.logger.Warn("peer health check failed", zap.String("peer", peer), zap.Error(err))
		}
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

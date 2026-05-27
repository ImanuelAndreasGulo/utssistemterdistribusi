package cluster

import (
	"net/http"
	"sync"
	"time"

	"distributed-counter-go/internal/config"
	"distributed-counter-go/internal/counter"

	"go.uber.org/zap"
)

// Manager coordinates cluster peers and synchronization.
type Manager struct {
	cfg     config.Config
	counter *counter.Counter
	logger  *zap.Logger
	client  *http.Client

	mu     sync.RWMutex
	peers  map[string]struct{}
	health map[string]bool
}

// NewManager creates a new cluster manager with initial peers.
func NewManager(cfg config.Config, ctr *counter.Counter, logger *zap.Logger) *Manager {
	m := &Manager{
		cfg:     cfg,
		counter: ctr,
		logger:  logger,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
		peers:  map[string]struct{}{},
		health: map[string]bool{},
	}

	for _, peer := range cfg.Peers {
		m.RegisterPeer(peer)
	}

	return m
}

// RegisterPeer adds a new peer to the cluster.
func (m *Manager) RegisterPeer(peer string) bool {
	if peer == "" || peer == m.cfg.NodeURL {
		return false
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.peers[peer]; exists {
		return false
	}

	m.peers[peer] = struct{}{}
	m.health[peer] = true
	return true
}

// GetPeers returns the list of peers.
func (m *Manager) GetPeers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	peers := make([]string, 0, len(m.peers))
	for peer := range m.peers {
		peers = append(peers, peer)
	}
	return peers
}

// GetHealthyPeers returns only peers that passed recent health checks.
func (m *Manager) GetHealthyPeers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	peers := make([]string, 0, len(m.peers))
	for peer := range m.peers {
		if m.health[peer] {
			peers = append(peers, peer)
		}
	}
	return peers
}

// SetHealth updates peer health status.
func (m *Manager) SetHealth(peer string, healthy bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.peers[peer]; ok {
		m.health[peer] = healthy
	}
}

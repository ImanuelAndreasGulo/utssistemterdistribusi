package config

import (
	"os"
	"strings"
	"time"
)

// Config holds runtime settings loaded from environment variables.
type Config struct {
	NodeID            string
	NodeURL           string
	Port              string
	DataPath          string
	Peers             []string
	HeartbeatInterval time.Duration
	SyncInterval      time.Duration
}

// Load reads configuration from environment variables with safe defaults.
func Load() Config {
	cfg := Config{
		NodeID:            getEnv("NODE_ID", "node1"),
		NodeURL:           getEnv("NODE_URL", "http://localhost:8001"),
		Port:              getEnv("PORT", "8001"),
		DataPath:          getEnv("DATA_PATH", "./data/node1.json"),
		Peers:             parsePeers(getEnv("PEERS", "")),
		HeartbeatInterval: getDuration("HEARTBEAT_INTERVAL", 5*time.Second),
		SyncInterval:      getDuration("SYNC_INTERVAL", 10*time.Second),
	}

	return cfg
}

func getEnv(key, def string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	return val
}

func parsePeers(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	peers := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			peers = append(peers, p)
		}
	}
	return peers
}

func getDuration(key string, def time.Duration) time.Duration {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return def
	}
	return d
}

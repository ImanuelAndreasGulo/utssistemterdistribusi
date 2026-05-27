package model

// IncrementRequest is the payload to increment local counter.
type IncrementRequest struct {
	Amount int `json:"amount"`
}

// IncrementResponse is the response for increment endpoint.
type IncrementResponse struct {
	Success      bool `json:"success"`
	Node         string `json:"node"`
	LocalCounter int `json:"local_counter"`
}

// CounterValueResponse returns global counter and active nodes.
type CounterValueResponse struct {
	GlobalCounter int `json:"global_counter"`
	ActiveNodes   int `json:"active_nodes"`
}

// JoinRequest is the payload to register a peer.
type JoinRequest struct {
	Peer string `json:"peer"`
}

// JoinResponse indicates result of join.
type JoinResponse struct {
	Message string `json:"message"`
}

// PeersResponse returns peer list.
type PeersResponse struct {
	Node  string   `json:"node"`
	Peers []string `json:"peers"`
}

// SyncRequest carries CRDT state to merge.
type SyncRequest struct {
	State map[string]int `json:"state"`
	Peer  string         `json:"peer,omitempty"`
}

// SyncResponse returns current state.
type SyncResponse struct {
	Node  string         `json:"node"`
	State map[string]int `json:"state"`
}

// HealthResponse returns health status.
type HealthResponse struct {
	Status string `json:"status"`
	Node   string `json:"node"`
}

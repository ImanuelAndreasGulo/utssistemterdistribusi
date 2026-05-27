# Distributed Counter API

## Endpoints

### POST /counter/increment
Increment local counter and broadcast update asynchronously.

**Request**
```json
{ "amount": 1 }
```

**Response**
```json
{ "success": true, "node": "node1", "local_counter": 15 }
```

### GET /counter/value
Aggregate global counter from peers.

**Response**
```json
{ "global_counter": 120, "active_nodes": 3 }
```

### POST /cluster/join
Register a new peer dynamically.

**Request**
```json
{ "peer": "http://localhost:8002" }
```

**Response**
```json
{ "message": "peer added" }
```

### GET /cluster/peers
List known peers.

**Response**
```json
{ "node": "node1", "peers": ["http://localhost:8002"] }
```

### GET /health
Health check endpoint.

**Response**
```json
{ "status": "ok", "node": "node1" }
```

### POST /sync
CRDT sync endpoint for merging state.

**Request**
```json
{ "peer": "http://localhost:8001", "state": {"node1": 5, "node2": 7} }
```

**Response**
```json
{ "node": "node2", "state": {"node1": 5, "node2": 7} }
```

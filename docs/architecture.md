# Architecture Overview

## Goals
- Distributed counter with eventual consistency
- No single point of failure
- Fault-tolerant peer communication

## Core Components
- **Counter (CRDT G-Counter)**: per-node counters, merge by max value per node
- **Cluster Manager**: peer discovery, health checks, sync and aggregation
- **Storage**: JSON file per node for persistence
- **API**: Gin-based REST endpoints

## Data Model
Each node stores a map of node IDs to integers:

```json
{ "node1": 5, "node2": 10 }
```

Global total = sum of all node values.

## Consistency Model
- Eventual consistency through periodic sync and broadcast on increment
- Conflict resolution via CRDT merge (max per node)

## Fault Tolerance
- Health checks with 3s timeout and retries
- Aggregation skips unreachable peers

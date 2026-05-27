# Distributed Counter System (Go)

Sistem counter terdistribusi berbasis Go dengan arsitektur multi-node REST API, menggunakan CRDT G-Counter untuk conflict resolution dan eventual consistency.

## Fitur
- Increment counter lokal dan broadcast async
- Agregasi counter global dari peer
- Registrasi peer dinamis
- Health check dan fault tolerance
- Penyimpanan JSON lokal per node

## Prasyarat
- Go 1.24+

## Konfigurasi
Gunakan file [.env](.env) sebagai contoh. Variabel yang digunakan:
- `NODE_ID`
- `NODE_URL`
- `PORT`
- `DATA_PATH`
- `PEERS`
- `HEARTBEAT_INTERVAL`
- `SYNC_INTERVAL`

## Menjalankan 2 Node (contoh)

**Node 1**
```bash
set NODE_ID=node1
set NODE_URL=http://localhost:8001
set PORT=8001
set DATA_PATH=./data/node1.json
set PEERS=http://localhost:8002
set HEARTBEAT_INTERVAL=5s
set SYNC_INTERVAL=10s

go run ./cmd/server
```

**Node 2**
```bash
set NODE_ID=node2
set NODE_URL=http://localhost:8002
set PORT=8002
set DATA_PATH=./data/node2.json
set PEERS=http://localhost:8001
set HEARTBEAT_INTERVAL=5s
set SYNC_INTERVAL=10s

go run ./cmd/server
```

## Endpoint
Detail endpoint tersedia di [docs/API.md](docs/API.md).

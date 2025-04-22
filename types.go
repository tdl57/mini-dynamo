// types.go
package main

import "time"

// Node represents a physical server in the cluster.
type Node struct {
	ID      string // unique identifier, e.g. "node‑1"
	Address string // host:port for gRPC comms
}

// VirtualNode is a point on the hash ring that maps back to a physical Node.
type VirtualNode struct {
	ID             string // e.g. "node‑1‑42"
	PhysicalNodeID string // the Node.ID it belongs to
	Position       uint64 // hash value around the ring
}

// Storage defines the interface for key-value persistence engines.
type Storage interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Close() error
}

type StorageNode struct {
	Node    Node
	Storage Storage
}

// VectorClock is a simple map-based vector clock for conflict tracking.
type VectorClock map[string]uint64

// KeyValuePair holds one record along with its causal metadata.
type KeyValuePair struct {
	Key       string
	Value     []byte
	Clock     VectorClock
	Timestamp time.Time
}

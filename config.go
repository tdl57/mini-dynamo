// config.go
package main

import "time"

// Config holds all tunable parameters for your Dynamo-style store.
type Config struct {
    ReplicationFactor           int           // how many nodes to replicate each key
    ReadQuorum                  int           // min. replicas that must agree on a read
    WriteQuorum                 int           // min. replicas that must confirm a write
    HeartbeatInterval           time.Duration // how often to ping peers
    VirtualNodesPerPhysicalNode int           // number of vnodes per physical Node
}

// DefaultConfig returns a sensible set of defaults:
//   - RF = 3 → quorum = 2
//   - heartbeat every 500ms
//   - 256 virtual nodes per physical node
func DefaultConfig() *Config {
    rf := 3
    q  := rf/2 + 1
    return &Config{
        ReplicationFactor:           rf,
        ReadQuorum:                  q,
        WriteQuorum:                 q,
        HeartbeatInterval:           500 * time.Millisecond,
        VirtualNodesPerPhysicalNode: 256,
    }
}
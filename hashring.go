// hash_ring.go
package main

import (
	"fmt"
	"hash/fnv"
	"sort"
)

type HashRing struct {
	virtualNodes       map[uint64]*VirtualNode // Map hash positions to virtual nodes
	sortedHashes       []uint64                // Sorted list of hash positions
	nodeToVirtualNodes map[string][]uint64     // Map node IDs to their virtual node hashes
}

// calculateHash returns a 64-bit FNV-1a hash of the input string.
func calculateHash(key string) uint64 {
	hasher := fnv.New64a()
	hasher.Write([]byte(key))
	return hasher.Sum64()
}

func NewHashRing() *HashRing {
	return &HashRing{
		virtualNodes:       make(map[uint64]*VirtualNode),
		sortedHashes:       []uint64{},
		nodeToVirtualNodes: make(map[string][]uint64),
	}
}

func (r *HashRing) AddNode(node *Node, numVirtualNodes int) {
	// Create virtual nodes and distribute them around the ring
	for i := 0; i < numVirtualNodes; i++ {
		vnodeID := fmt.Sprintf("%s-%d", node.ID, i)
		hash := calculateHash(vnodeID) // Using consistent hash function like murmur3

		vnode := &VirtualNode{
			ID:             vnodeID,
			PhysicalNodeID: node.ID,
			Position:       hash,
		}

		r.virtualNodes[hash] = vnode
		r.nodeToVirtualNodes[node.ID] = append(r.nodeToVirtualNodes[node.ID], hash)
	}

	// Update sorted hashes
	r.sortedHashes = make([]uint64, 0, len(r.virtualNodes))
	for hash := range r.virtualNodes {
		r.sortedHashes = append(r.sortedHashes, hash)
	}
	sort.Slice(r.sortedHashes, func(i, j int) bool {
		return r.sortedHashes[i] < r.sortedHashes[j]
	})
}

func (r *HashRing) GetNodesForKey(key string, n int) []*VirtualNode {
	if len(r.sortedHashes) == 0 {
		return nil
	}
	hash := calculateHash(key)
	index := sort.Search(len(r.sortedHashes), func(i int) bool {
		return r.sortedHashes[i] >= hash
	})
	if index == len(r.sortedHashes) {
		index = 0
	}

	result := make([]*VirtualNode, 0, n)
	distinctNodes := make(map[string]bool)

	for len(result) < n && len(distinctNodes) < len(r.nodeToVirtualNodes) {
		vnode := r.virtualNodes[r.sortedHashes[index]]
		if !distinctNodes[vnode.PhysicalNodeID] {
			distinctNodes[vnode.PhysicalNodeID] = true
			result = append(result, vnode)
		}
		index = (index + 1) % len(r.sortedHashes)
	}
	return result
}

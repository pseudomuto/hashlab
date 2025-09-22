package hashring

import (
	"errors"
	"fmt"
	"hash/crc32"
	"slices"
	"sort"
	"sync"
)

// HashRing represents a consistent hash ring
type HashRing struct {
	mu         sync.RWMutex
	ring       map[uint32]string // hash -> server name
	serverKeys []uint32          // sorted server hashes
	servers    map[string]bool   // map of servers
}

// New creates a new hash ring with the specified number of virtual nodes per server
func New(virtualNodes int) *HashRing {
	return &HashRing{
		ring:       make(map[uint32]string),
		serverKeys: make([]uint32, 0),
		servers:    make(map[string]bool),
		// TODO: Keep track of virtual nodes
	}
}

// hashKey generates a hash value for the given key
func (h *HashRing) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// AddServer adds a server to the hash ring
func (h *HashRing) AddServer(server string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.servers[server] {
		return fmt.Errorf("server already exists: %s", server)
	}

	h.servers[server] = true

	// TODO: Replace this implementation (following 3 lines).
	//
	// We want to create entries in ring and serverKeys for each virtual node
	// vnodes can be named however you like, my suggestion would be `<server>#<n>`, where server is the parameter and N
	// is the virtual node index,
	hash := h.hashKey(server)
	h.ring[hash] = server
	h.serverKeys = append(h.serverKeys, hash)

	slices.Sort(h.serverKeys)
	return nil
}

// RemoveServer removes a server from the hash ring
func (h *HashRing) RemoveServer(server string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.servers[server] {
		return fmt.Errorf("server not found: %s", server)
	}

	delete(h.servers, server)

	// TODO: Replace the code between BEGIN and END comments.
	//
	// We need to remove all vnodes from ring and serverKeys for the specified server. Be sure to maintain sort order.
	// BEGIN
	hash := h.hashKey(server)
	delete(h.ring, hash)

	idx := slices.Index(h.serverKeys, hash)
	h.serverKeys = append(h.serverKeys[:idx], h.serverKeys[idx+1:]...)
	// END

	return nil
}

// GetServer returns the server responsible for the given key
func (h *HashRing) GetServer(key string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.servers) == 0 {
		return "", errors.New("no servers available")
	}

	hash := h.hashKey(key)

	// NB: This works because we keep serverKeys sorted.
	idx := sort.Search(len(h.serverKeys), func(i int) bool {
		return h.serverKeys[i] >= hash
	})

	// NB: sort.Search returns n, when not found. Because we know a server exists, it must be at index 0.
	if idx == len(h.serverKeys) {
		idx = 0 // Ringify
	}

	return h.ring[h.serverKeys[idx]], nil
}

// GetServers returns all servers in the ring
func (h *HashRing) GetServers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	servers := make([]string, 0, len(h.servers))
	for server := range h.servers {
		servers = append(servers, server)
	}

	slices.Sort(servers)
	return servers
}

// GetDistribution returns a map of server -> count of keys
// This is useful for analyzing how well keys are distributed
func (h *HashRing) GetDistribution(keys []string) map[string]int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	dist := make(map[string]int, len(h.servers))
	for server := range h.servers {
		dist[server] = 0
	}

	for _, key := range keys {
		server, err := h.GetServer(key)
		if err == nil {
			dist[server]++
		}
	}

	return dist
}

// Size returns the number of physical servers in the ring
func (h *HashRing) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.servers)
}

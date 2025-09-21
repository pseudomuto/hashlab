package hashring

import (
	"errors"
	"fmt"
	"hash/crc32"
	"slices"
	"sync"
)

// HashRing represents a consistent hash ring
type HashRing struct {
	mu sync.RWMutex
	// TODO: Update to use ring (map[uint32]string), serverKeys ([]uint32), and servers (same).
	servers    map[string]bool
	serverList []string
}

// New creates a new hash ring with the specified number of virtual nodes per server
func New(virtualNodes int) *HashRing {
	return &HashRing{
		// TODO: Update struct fields
		servers:    make(map[string]bool),
		serverList: make([]string, 0),
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

	// TODO: Remove the next two lines (field no longer exists).
	h.serverList = append(h.serverList, server)
	slices.Sort(h.serverList) // For consistent indexing

	// TODO: Hash the server name directly.
	// TODO: Add the hash to the ring.
	// TODO: Add the key to serverKeys.
	// TODO: Ensure serverKeys remains in sorted order.

	return nil
}

// RemoveServer removes a server from the hash ring
func (h *HashRing) RemoveServer(server string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.servers[server]; !ok {
		return fmt.Errorf("server not found: %s", server)
	}

	delete(h.servers, server)

	////////////////////////////////////////////////////////////////////////////////
	// TODO: ❌ THIS IS THE PROBLEM: hash % N
	// When N changes (add/remove server), most keys get different results!
	//
	// REMOVE THESE TWO LINES
	////////////////////////////////////////////////////////////////////////////////
	idx := slices.Index(h.serverList, server)
	h.serverList = append(h.serverList[:idx], h.serverList[idx+1:]...)

	// TODO: Hash the server
	// TODO: Remove it from the ring map
	// TODO: Remove it from the serverKeys list (retain sort order)

	return nil
}

// GetServer returns the server responsible for the given key
func (h *HashRing) GetServer(key string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.serverList) == 0 {
		return "", errors.New("no servers available")
	}

	hash := h.hashKey(key)

	// TODO: Remove the next two lines
	idx := int(hash % uint32(len(h.serverList))) // nolint: gosec
	return h.serverList[idx], nil

	// TODO: Find the index of the next ("clockwise") server in the ring (see sort.Search).
	// TODO: Any special case you can think of?
	// TODO: Return the server from the ring (you need the preceeding values).
}

// GetServers returns all servers in the ring
func (h *HashRing) GetServers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// TODO: serverList is no more, remove the next 3 lines.
	servers := make([]string, len(h.serverList))
	copy(servers, h.serverList)
	return servers

	// TODO: Return a sorted list of server names.
}

// GetDistribution returns a map of server -> count of keys
// This is useful for analyzing how well keys are distributed
func (h *HashRing) GetDistribution(keys []string) map[string]int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	dist := make(map[string]int, len(h.serverList))
	for _, server := range h.serverList {
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

	return len(h.serverList)
}

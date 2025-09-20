package hashring

import (
	"errors"
	"fmt"
	"hash/crc32"
	"slices"
	"sort"
	"sync"
	"time"
)

// HashRing represents a consistent hash ring for distributed systems.
// It uses virtual nodes to ensure even distribution of keys across servers
// and maintains consistency when servers are added or removed.
//
// The ring is thread-safe and supports concurrent operations.
type HashRing struct {
	mu         sync.RWMutex
	ring       map[uint32]string // hash position -> server name
	serverKeys []uint32          // sorted hash positions
	servers    map[string]bool   // set of server names
	vnodes     int               // number of virtual nodes per server
}

// New creates a new hash ring with the specified number of virtual nodes per server.
//
// The virtualNodes parameter determines how many positions each physical server
// will occupy on the hash ring. More virtual nodes provide better distribution
// but use more memory. Recommended values:
//   - Small clusters (3-10 servers): 100-200 virtual nodes
//   - Medium clusters (10-50 servers): 50-150 virtual nodes
//   - Large clusters (50+ servers): 20-100 virtual nodes
//
// Example:
//
//	ring := hashring.New(150)
//	ring.AddServer("server1")
//	ring.AddServer("server2")
//	server, _ := ring.GetServer("mykey")
func New(virtualNodes int) *HashRing {
	return &HashRing{
		ring:       make(map[uint32]string),
		serverKeys: make([]uint32, 0),
		servers:    make(map[string]bool),
		vnodes:     virtualNodes,
	}
}

// hashKey generates a hash value for the given key
func (h *HashRing) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// AddServer adds a server to the hash ring.
//
// The server is distributed across multiple positions on the ring using virtual nodes.
// This operation is thread-safe and will update the sorted key list for efficient lookups.
//
// Returns an error if the server already exists in the ring.
//
// Example:
//
//	err := ring.AddServer("cache-server-1")
//	if err != nil {
//		log.Printf("Failed to add server: %v", err)
//	}
func (h *HashRing) AddServer(server string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.servers[server] {
		return fmt.Errorf("server %s already exists", server)
	}

	h.servers[server] = true

	// Add virtual nodes for this server
	for i := 0; i < h.vnodes; i++ {
		hash := h.hashKey(fmt.Sprintf("%s#%d", server, i))
		h.ring[hash] = server
		h.serverKeys = append(h.serverKeys, hash)
	}

	// Sort the keys
	slices.Sort(h.serverKeys)
	return nil
}

// RemoveServer removes a server from the hash ring.
//
// All virtual nodes associated with the server are removed, and keys previously
// mapped to this server will be redistributed to the remaining servers.
// This operation is thread-safe.
//
// Returns an error if the server does not exist in the ring.
//
// Example:
//
//	err := ring.RemoveServer("cache-server-1")
//	if err != nil {
//		log.Printf("Failed to remove server: %v", err)
//	}
func (h *HashRing) RemoveServer(server string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.servers[server] {
		return fmt.Errorf("server %s does not exist", server)
	}

	delete(h.servers, server)

	for i := range h.vnodes {
		hash := h.hashKey(fmt.Sprintf("%s#%d", server, i))
		delete(h.ring, hash)

		idx := slices.Index(h.serverKeys, hash)
		h.serverKeys = append(h.serverKeys[:idx], h.serverKeys[idx+1:]...)
	}

	return nil
}

// GetServer returns the server responsible for the given key.
//
// Uses consistent hashing to determine which server should handle the key.
// The same key will always map to the same server (unless the ring changes).
// This operation is thread-safe and uses binary search for O(log n) lookup time.
//
// Returns an error if the hash ring is empty.
//
// Example:
//
//	server, err := ring.GetServer("user:12345")
//	if err != nil {
//		log.Printf("No servers available: %v", err)
//		return
//	}
//	fmt.Printf("Key 'user:12345' maps to %s\n", server)
func (h *HashRing) GetServer(key string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.ring) == 0 {
		return "", errors.New("hash ring is empty")
	}

	hash := h.hashKey(key)

	// Binary search to find the first server clockwise from the key's hash
	idx := sort.Search(len(h.serverKeys), func(i int) bool {
		return h.serverKeys[i] >= hash
	})

	// Wrap around if we've gone past the end
	if idx == len(h.serverKeys) {
		idx = 0
	}

	return h.ring[h.serverKeys[idx]], nil
}

// GetServers returns a sorted list of all servers currently in the ring.
//
// This operation is thread-safe and returns a new slice to prevent external
// modification of the internal server list.
//
// Example:
//
//	servers := ring.GetServers()
//	fmt.Printf("Active servers: %v\n", servers)
func (h *HashRing) GetServers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	servers := make([]string, 0, len(h.servers))
	for server := range h.servers {
		servers = append(servers, server)
	}

	sort.Strings(servers)
	return servers
}

// GetDistribution analyzes how a set of keys would be distributed across servers.
//
// Returns a map where each key is a server name and the value is the count of
// keys that would be assigned to that server. This is useful for:
//   - Analyzing load balance quality
//   - Capacity planning
//   - Debugging hot spots
//
// This operation is thread-safe but may be slow for large key sets.
//
// Example:
//
//	keys := []string{"user:1", "user:2", "user:3", "session:abc"}
//	dist := ring.GetDistribution(keys)
//	for server, count := range dist {
//		fmt.Printf("%s: %d keys\n", server, count)
//	}
func (h *HashRing) GetDistribution(keys []string) map[string]int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	distribution := make(map[string]int)
	for server := range h.servers {
		distribution[server] = 0
	}

	for _, key := range keys {
		server, err := h.GetServer(key)
		if err == nil {
			distribution[server]++
		}
	}

	return distribution
}

// Size returns the number of physical servers in the ring.
//
// This counts actual servers, not virtual nodes. For the total number of
// virtual nodes, multiply Size() by the virtualNodes parameter used in New().
// This operation is thread-safe.
//
// Example:
//
//	if ring.Size() == 0 {
//		fmt.Println("No servers available")
//	}
func (h *HashRing) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.servers)
}

// AnalyzePerformance runs a comprehensive performance analysis on the hash ring.
//
// This method evaluates:
//   - Key distribution across servers (uniformity)
//   - Average lookup latency per key
//   - Distribution quality using Coefficient of Variation (CV)
//
// A lower CV percentage indicates better distribution:
//   - CV < 5%: Excellent distribution
//   - CV < 10%: Good distribution
//   - CV > 10%: Consider adjusting virtual nodes
//
// This operation is thread-safe but may be slow for large key sets.
// It's recommended to run this during testing or monitoring, not in hot paths.
//
// Example:
//
//	testKeys := generateTestKeys(10000)
//	metrics := ring.AnalyzePerformance(testKeys)
//	metrics.Print() // Display formatted analysis
func (h *HashRing) AnalyzePerformance(keys []string) PerformanceMetrics {
	start := time.Now()

	// Measure average latency
	distribution := h.GetDistribution(keys)
	avgLatency := time.Since(start) / time.Duration(len(keys))

	// Calculate distribution quality (Coefficient of Variation)
	mean := float64(len(keys)) / float64(len(distribution))
	var variance float64
	for _, count := range distribution {
		diff := float64(count) - mean
		variance += diff * diff
	}

	stdDev := 0.0
	if len(distribution) > 0 {
		variance /= float64(len(distribution))
		stdDev = variance
		for range 10 { // Simple sqrt approximation
			stdDev = (stdDev + variance/stdDev) / 2
		}
	}

	cv := 0.0
	if mean > 0 {
		cv = (stdDev / mean) * 100
	}

	return PerformanceMetrics{
		TotalKeys:      len(keys),
		Servers:        len(distribution),
		AvgLatency:     avgLatency,
		DistributionCV: cv,
		Distribution:   distribution,
	}
}

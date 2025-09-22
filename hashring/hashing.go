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

// HashRing represents a consistent hash ring
type HashRing struct {
	mu         sync.RWMutex
	ring       map[uint32]string // hash -> server name
	serverKeys []uint32          // sorted server hashes
	servers    map[string]bool   // map of servers
	vnodes     int               // The number of virtual nodes per server
}

// New creates a new hash ring with the specified number of virtual nodes per server
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

// AddServer adds a server to the hash ring
func (h *HashRing) AddServer(server string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.servers[server] {
		return fmt.Errorf("server already exists: %s", server)
	}

	h.servers[server] = true

	for i := range h.vnodes {
		hash := h.hashKey(fmt.Sprintf("%s#%d", server, i))
		h.ring[hash] = server
		h.serverKeys = append(h.serverKeys, hash)
	}

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

	for i := range h.vnodes {
		hash := h.hashKey(fmt.Sprintf("%s#%d", server, i))
		delete(h.ring, hash)

		idx := slices.Index(h.serverKeys, hash)
		h.serverKeys = append(h.serverKeys[:idx], h.serverKeys[idx+1:]...)
	}

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
// See: https://en.wikipedia.org/wiki/Coefficient_of_variation
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

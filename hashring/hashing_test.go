package hashring

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	ring := New(150)
	require.NotNil(t, ring, "New returned nil")
	require.Equal(t, 150, ring.vnodes, "Expected 150 virtual nodes")
}

func TestAddServer(t *testing.T) {
	ring := New(150)

	err := ring.AddServer("server1")
	require.NoError(t, err, "Failed to add server")
	require.Equal(t, 1, ring.Size(), "Expected 1 server")

	// Adding the same server should fail
	err = ring.AddServer("server1")
	require.Error(t, err, "Expected error when adding duplicate server")

	// Add more servers
	require.NoError(t, ring.AddServer("server2"))
	require.NoError(t, ring.AddServer("server3"))
	require.Equal(t, 3, ring.Size(), "Expected 3 servers")
}

func TestRemoveServer(t *testing.T) {
	ring := New(150)
	require.NoError(t, ring.AddServer("server1"))
	require.NoError(t, ring.AddServer("server2"))
	require.NoError(t, ring.AddServer("server3"))

	err := ring.RemoveServer("server2")
	require.NoError(t, err, "Failed to remove server")
	require.Equal(t, 2, ring.Size(), "Expected 2 servers")

	// Removing non-existent server should fail
	err = ring.RemoveServer("server2")
	require.Error(t, err, "Expected error when removing non-existent server")
}

func TestGetServer(t *testing.T) {
	ring := New(150)

	// Empty ring should return error
	_, err := ring.GetServer("key1")
	require.Error(t, err, "Expected error for empty ring")

	require.NoError(t, ring.AddServer("server1"))
	require.NoError(t, ring.AddServer("server2"))
	require.NoError(t, ring.AddServer("server3"))

	// Test that the same key always maps to the same server
	key := "test-key"
	server1, err := ring.GetServer(key)
	require.NoError(t, err, "Failed to get server")

	server2, err := ring.GetServer(key)
	require.NoError(t, err, "Failed to get server")

	require.Equal(t, server1, server2, "Same key mapped to different servers")
}

func TestConsistency(t *testing.T) {
	ring := New(150)
	require.NoError(t, ring.AddServer("server1"))
	require.NoError(t, ring.AddServer("server2"))
	require.NoError(t, ring.AddServer("server3"))

	// Map 1000 keys
	keyToServer := make(map[string]string)
	for i := range 1000 {
		key := fmt.Sprintf("key-%d", i)
		server, err := ring.GetServer(key)
		require.NoError(t, err)
		keyToServer[key] = server
	}

	// Add a new server
	require.NoError(t, ring.AddServer("server4"))

	// Check how many keys moved
	moved := 0
	for key, oldServer := range keyToServer {
		newServer, err := ring.GetServer(key)
		require.NoError(t, err)
		if newServer != oldServer {
			moved++
		}
	}

	// With consistent hashing, roughly 1/4 of keys should move (1000 / 4 = 250)
	// Allow some variance (150-350)
	require.GreaterOrEqual(t, moved, 150, "Too few keys moved")
	require.LessOrEqual(t, moved, 350, "Too many keys moved")

	t.Logf("Keys moved when adding server: %d out of 1000 (%.1f%%)", moved, float64(moved)/10)
}

func TestDistribution(t *testing.T) {
	ring := New(150)
	require.NoError(t, ring.AddServer("server1"))
	require.NoError(t, ring.AddServer("server2"))
	require.NoError(t, ring.AddServer("server3"))

	// Generate test keys
	keys := make([]string, 10000)
	for i := range 10000 {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	distribution := ring.GetDistribution(keys)

	// Check that distribution is relatively even
	expectedPerServer := 10000 / 3
	for server, count := range distribution {
		// Allow 20% variance
		minExpected := int(float64(expectedPerServer) * 0.8)
		maxExpected := int(float64(expectedPerServer) * 1.2)

		require.GreaterOrEqual(t, count, minExpected,
			"Server %s has too few keys: %d (expected >=%d)", server, count, minExpected)
		require.LessOrEqual(t, count, maxExpected,
			"Server %s has too many keys: %d (expected <=%d)", server, count, maxExpected)
	}

	t.Logf("Distribution: %v", distribution)
}

func TestDistributionStandardDeviation(t *testing.T) {
	ring := New(150)
	require.NoError(t, ring.AddServer("server1"))
	require.NoError(t, ring.AddServer("server2"))
	require.NoError(t, ring.AddServer("server3"))
	require.NoError(t, ring.AddServer("server4"))

	// Generate test keys
	keys := make([]string, 10000)
	for i := range 10000 {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	distribution := ring.GetDistribution(keys)

	// Calculate mean
	var sum int
	for _, count := range distribution {
		sum += count
	}
	mean := float64(sum) / float64(len(distribution))

	// Calculate standard deviation
	var variance float64
	for _, count := range distribution {
		diff := float64(count) - mean
		variance += diff * diff
	}
	variance /= float64(len(distribution))
	stdDev := math.Sqrt(variance)

	// Standard deviation should be reasonable (< 20% of mean for small server count)
	require.LessOrEqual(t, stdDev, mean*0.2,
		"Standard deviation too high: %.2f (mean: %.2f)", stdDev, mean)

	t.Logf("Distribution stats - Mean: %.2f, StdDev: %.2f", mean, stdDev)
}

func TestVirtualNodesImpact(t *testing.T) {
	keys := make([]string, 10000)
	for i := range 10000 {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	testCases := []int{10, 50, 150, 500}

	for _, vnodes := range testCases {
		ring := New(vnodes)
		require.NoError(t, ring.AddServer("server1"))
		require.NoError(t, ring.AddServer("server2"))
		require.NoError(t, ring.AddServer("server3"))

		distribution := ring.GetDistribution(keys)

		// Calculate coefficient of variation
		var sum, mean, variance float64
		for _, count := range distribution {
			sum += float64(count)
		}
		mean = sum / float64(len(distribution))

		for _, count := range distribution {
			diff := float64(count) - mean
			variance += diff * diff
		}
		stdDev := math.Sqrt(variance / float64(len(distribution)))
		cv := (stdDev / mean) * 100 // coefficient of variation as percentage

		t.Logf("Virtual nodes: %d, CV: %.2f%%", vnodes, cv)

		// More virtual nodes should result in better distribution (lower CV)
		if vnodes >= 150 {
			require.LessOrEqual(t, cv, 10.0,
				"Coefficient of variation too high for %d vnodes: %.2f%%", vnodes, cv)
		}
	}
}

func TestConcurrency(t *testing.T) {
	ring := New(150)
	require.NoError(t, ring.AddServer("server1"))
	require.NoError(t, ring.AddServer("server2"))

	// Test concurrent reads
	done := make(chan bool)
	for i := range 10 {
		go func(id int) {
			for j := range 100 {
				key := fmt.Sprintf("key-%d-%d", id, j)
				_, err := ring.GetServer(key)
				require.NoError(t, err, "Error getting server")
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for range 10 {
		<-done
	}
}

package main

import (
	"fmt"

	"github.com/pseudomuto/hashlab/hashring"
)

// Example 1: Simple cache distribution
func main() {
	fmt.Println("=== Cache Distribution Example ===")
	fmt.Println()

	ring := hashring.New(150)

	// Add cache servers
	cacheServers := []string{"cache-1", "cache-2", "cache-3"}
	for _, server := range cacheServers {
		ring.AddServer(server)
	}

	// Simulate cache key lookups
	cacheKeys := []string{
		"user:1234:profile",
		"user:5678:preferences",
		"session:abc123",
		"product:9999",
		"cart:user-42",
	}

	fmt.Println("Cache key routing:")
	for _, key := range cacheKeys {
		server, _ := ring.GetServer(key)
		fmt.Printf("  %s → %s\n", key, server)
	}

	// Consistent mapping - same key always goes to same server
	fmt.Println("\nVerifying consistency:")
	for i := range 3 {
		server, _ := ring.GetServer("user:1234:profile")
		fmt.Printf("  Attempt %d: user:1234:profile → %s\n", i+1, server)
	}
}

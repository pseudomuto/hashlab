package main

import (
	"fmt"
	"log"
	"math"

	"github.com/pseudomuto/hashlab/hashring"
)

const (
	numKeys   = 10_000
	numVNodes = 150
)

func main() {
	fmt.Println("=== Consistent Hash Ring Demo ===")

	// Create a hash ring with 150 virtual nodes per server
	ring := hashring.New(numVNodes)

	// Add servers
	fmt.Println("Adding servers...")
	servers := []string{"server-A", "server-B", "server-C"}
	for _, server := range servers {
		if err := ring.AddServer(server); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  ✓ Added %s\n", server)
	}

	// Generate sample keys
	fmt.Println("\nGenerating 10,000 sample keys...")
	keys := make([]string, numKeys)
	for i := range numKeys {
		keys[i] = fmt.Sprintf("user-%d", i)
	}

	// Show key distribution
	fmt.Println("\nKey distribution across servers:")
	distribution := ring.GetDistribution(keys)
	for server, count := range distribution {
		percentage := float64(count) / float64(len(keys)) * 100
		fmt.Printf("  %s: %d keys (%.2f%%)\n", server, count, percentage)
	}

	// Calculate distribution quality
	mean := float64(len(keys)) / float64(len(servers))
	var variance float64
	for _, count := range distribution {
		diff := float64(count) - mean
		variance += diff * diff
	}
	stdDev := math.Sqrt(variance / float64(len(servers)))
	cv := (stdDev / mean) * 100

	fmt.Printf("\nDistribution quality (Coefficient of Variation): %.2f%%\n", cv)
	if cv < 5 {
		fmt.Println("  ✓ Excellent distribution!")
	} else if cv < 10 {
		fmt.Println("  ✓ Good distribution")
	} else {
		fmt.Println("  ⚠ Consider using more virtual nodes")
	}

	// Map some specific keys
	fmt.Println("\nExample key mappings:")
	exampleKeys := []string{"user-42", "user-1337", "user-9999", "session-abc123"}
	for _, key := range exampleKeys {
		server, _ := ring.GetServer(key)
		fmt.Printf("  %s → %s\n", key, server)
	}

	// Demonstrate adding a new server
	fmt.Println("\n--- Adding a new server (server-D) ---")

	// Track which keys move
	oldMapping := make(map[string]string)
	for _, key := range keys {
		server, _ := ring.GetServer(key)
		oldMapping[key] = server
	}

	if err := ring.AddServer("server-D"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("  ✓ Added server-D")

	moved := 0
	for _, key := range keys {
		newServer, _ := ring.GetServer(key)
		if newServer != oldMapping[key] {
			moved++
		}
	}

	fmt.Printf("\nKeys that moved: %d out of %d (%.2f%%)\n", moved, len(keys), float64(moved)/float64(len(keys))*100)
	expectedMove := float64(len(keys)) / float64(len(servers)+1)
	fmt.Printf("Expected keys to move: ~%.0f (%.2f%%)\n", expectedMove, expectedMove/float64(len(keys))*100)

	// Show new distribution
	fmt.Println("\nNew key distribution:")
	distribution = ring.GetDistribution(keys)
	for server, count := range distribution {
		percentage := float64(count) / float64(len(keys)) * 100
		fmt.Printf("  %s: %d keys (%.2f%%)\n", server, count, percentage)
	}

	// Demonstrate removing a server
	fmt.Println("\n--- Removing a server (server-B) ---")

	// Track which keys move
	oldMapping = make(map[string]string)
	for _, key := range keys {
		server, _ := ring.GetServer(key)
		oldMapping[key] = server
	}

	if err := ring.RemoveServer("server-B"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("  ✓ Removed server-B")

	moved = 0
	movedTo := make(map[string]int)
	for _, key := range keys {
		newServer, _ := ring.GetServer(key)
		if newServer != oldMapping[key] {
			moved++
			movedTo[newServer]++
		}
	}

	fmt.Printf("\nKeys that moved: %d out of %d (%.2f%%)\n", moved, len(keys), float64(moved)/float64(len(keys))*100)
	fmt.Println("\nKeys moved to:")
	for server, count := range movedTo {
		fmt.Printf("  %s: %d keys\n", server, count)
	}

	// Show final distribution
	fmt.Println("\nFinal key distribution:")
	distribution = ring.GetDistribution(keys)
	for server, count := range distribution {
		percentage := float64(count) / float64(len(keys)) * 100
		fmt.Printf("  %s: %d keys (%.2f%%)\n", server, count, percentage)
	}

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("  • Virtual nodes ensure even distribution")
	fmt.Println("  • Adding servers only moves ~1/N keys")
	fmt.Println("  • Removing servers redistributes only affected keys")
	fmt.Println("  • Consistent hashing minimizes disruption during scaling")
}

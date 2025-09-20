package main

import (
	"fmt"

	"github.com/pseudomuto/hashlab/hashring"
)

// Example comparing hashring with naive modulo hashing
func main() {
	fmt.Println()
	fmt.Println("=== Consistent Hash vs Modulo Comparison ===")
	fmt.Println()

	numServers := 3
	numKeys := 100

	// Modulo hashing
	moduloMap := make(map[string]int)
	for i := range numKeys {
		key := fmt.Sprintf("key-%d", i)
		// Simulate simple modulo: hash(key) % numServers
		serverIdx := i % numServers
		moduloMap[key] = serverIdx
	}

	// Add a server with modulo
	newNumServers := 4
	moduloMoved := 0
	for i := range numKeys {
		key := fmt.Sprintf("key-%d", i)
		newServerIdx := i % newNumServers
		if newServerIdx != moduloMap[key] {
			moduloMoved++
		}
	}

	// Consistent hashing
	ring := hashring.New(150)
	for i := range numServers {
		ring.AddServer(fmt.Sprintf("server-%d", i))
	}

	chMap := make(map[string]string)
	for i := range numKeys {
		key := fmt.Sprintf("key-%d", i)
		server, _ := ring.GetServer(key)
		chMap[key] = server
	}

	// Add a server with consistent hashing
	ring.AddServer("server-3")
	chMoved := 0
	for i := range numKeys {
		key := fmt.Sprintf("key-%d", i)
		newServer, _ := ring.GetServer(key)
		if newServer != chMap[key] {
			chMoved++
		}
	}

	fmt.Printf("When adding 1 server to %d servers:\n", numServers)
	fmt.Printf("  Modulo hashing: %d keys moved (%.0f%%)\n",
		moduloMoved, float64(moduloMoved)/float64(numKeys)*100)
	fmt.Printf("  Consistent hashing: %d keys moved (%.0f%%)\n",
		chMoved, float64(chMoved)/float64(numKeys)*100)
	fmt.Printf("\nConsistent hashing moved %.1fx fewer keys!\n",
		float64(moduloMoved)/float64(chMoved))
}

package main

import (
	"fmt"

	"github.com/pseudomuto/hashlab/hashring"
)

// Example of a Load balancer with sticky sessions.
func main() {
	fmt.Println()
	fmt.Println("=== Load Balancer Example ===")
	fmt.Println()

	ring := hashring.New(150)

	// Add backend servers
	backends := []string{"backend-1:8080", "backend-2:8080", "backend-3:8080"}
	for _, backend := range backends {
		ring.AddServer(backend)
	}

	// Route requests based on client IP or session ID
	fmt.Println("Request routing (by session ID):")
	sessions := []string{
		"session-abc123",
		"session-def456",
		"session-ghi789",
		"session-jkl012",
		"session-mno345",
	}

	for _, sessionID := range sessions {
		backend, _ := ring.GetServer(sessionID)
		fmt.Printf("  %s → %s\n", sessionID, backend)
	}

	// Sticky sessions - same session always goes to same backend
	fmt.Println("\nVerifying sticky sessions (session-abc123):")
	for i := range 3 {
		backend, _ := ring.GetServer("session-abc123")
		fmt.Printf("  Request %d → %s\n", i+1, backend)
	}
}

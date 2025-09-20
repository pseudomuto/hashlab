package main

import (
	"fmt"

	"github.com/pseudomuto/hashlab/hashring"
)

// Example of Database sharding
func main() {
	fmt.Println()
	fmt.Println("=== Database Sharding Example ===")
	fmt.Println()

	ring := hashring.New(150)

	// Add database shards
	shards := []string{"db-shard-1", "db-shard-2", "db-shard-3", "db-shard-4"}
	for _, shard := range shards {
		ring.AddServer(shard)
	}

	// Route user data to shards
	fmt.Println("User data routing:")
	for userID := 1; userID <= 10; userID++ {
		key := fmt.Sprintf("user:%d", userID)
		shard, _ := ring.GetServer(key)
		fmt.Printf("  User %d → %s\n", userID, shard)
	}

	// Show what happens when we add a new shard
	fmt.Println("\nAdding a new shard (db-shard-5)...")

	beforeMap := make(map[string]string)
	for userID := 1; userID <= 100; userID++ {
		key := fmt.Sprintf("user:%d", userID)
		shard, _ := ring.GetServer(key)
		beforeMap[key] = shard
	}

	ring.AddServer("db-shard-5")

	moved := 0
	for userID := 1; userID <= 100; userID++ {
		key := fmt.Sprintf("user:%d", userID)
		newShard, _ := ring.GetServer(key)
		if newShard != beforeMap[key] {
			moved++
		}
	}

	fmt.Printf("Users that need to be migrated: %d out of 100 (%.0f%%)\n",
		moved, float64(moved))
}

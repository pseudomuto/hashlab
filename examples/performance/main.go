package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/pseudomuto/hashlab/hashring"
)

var (
	numKeys    int
	numServers int
	numVNodes  int
)

func init() {
	flag.IntVar(&numKeys, "keys", 10_000, "The number of keys to test")
	flag.IntVar(&numServers, "servers", 3, "The number of servers to distribute across")
	flag.IntVar(&numVNodes, "vnodes", 150, "The number of virtual nodes per server")
	flag.Parse()
}

func main() {
	// Create hash ring
	ring := hashring.New(numKeys)

	// Add servers
	for i := range numServers {
		if err := ring.AddServer(fmt.Sprintf("server-%d", i+1)); err != nil {
			log.Fatal(err)
		}
	}

	// Generate test keys
	keys := make([]string, numKeys)
	for i := range numKeys {
		keys[i] = fmt.Sprintf("user-%d", i+1)
	}

	// Analyze performance
	ring.AnalyzePerformance(keys).Print()
}

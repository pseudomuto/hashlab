package hashring

import (
	"fmt"
	"time"
)

// PerformanceMetrics contains performance analysis results
//
// === Optimization Tips ===
//
// 1. Virtual Nodes:
//   - Too few (<50): Poor distribution
//   - Optimal (100-500): Good balance
//   - Too many (>1000): Memory overhead
//
// 2. For Better Performance:
//   - Use RWMutex for read-heavy loads
//   - Pre-allocate slices with capacity
//   - Consider caching hash values
//
// 3. Production Considerations:
//   - Monitor distribution quality (CV < 10%)
//   - Set timeouts for operations
//   - Plan gradual migration strategies
//   - Test with production-like data
type PerformanceMetrics struct {
	TotalKeys      int
	Servers        int
	VirtualNodes   int
	AvgLatency     time.Duration
	DistributionCV float64 // Coefficient of Variation
	Distribution   map[string]int
}

// Print displays a formatted performance analysis report to stdout.
//
// The report includes:
//   - Total number of keys analyzed
//   - Number of servers in the ring
//   - Average latency per key lookup
//   - Distribution quality (Coefficient of Variation)
//   - Per-server key distribution with percentages
//
// The distribution quality is evaluated as:
//   - CV < 5%: Excellent distribution (✅)
//   - CV < 10%: Good distribution (✅)
//   - CV >= 10%: Poor distribution, consider more virtual nodes (⚠️)
//
// Example output:
//
//	=== Performance Analysis ===
//	Total Keys: 10000
//	Servers: 3
//	Avg Latency: 125ns per key
//	Distribution CV: 3.45%
//	✅ Excellent distribution!
//
//	Key Distribution:
//	  server-1: 3342 keys (33.4%)
//	  server-2: 3321 keys (33.2%)
//	  server-3: 3337 keys (33.4%)
func (metrics PerformanceMetrics) Print() {
	fmt.Println("\n=== Performance Analysis ===")
	fmt.Printf("Total Keys: %d\n", metrics.TotalKeys)
	fmt.Printf("Servers: %d\n", metrics.Servers)
	fmt.Printf("Avg Latency: %v per key\n", metrics.AvgLatency)
	fmt.Printf("Distribution CV: %.2f%%\n", metrics.DistributionCV)

	if metrics.DistributionCV < 5 {
		fmt.Println("✅ Excellent distribution!")
	} else if metrics.DistributionCV < 10 {
		fmt.Println("✅ Good distribution")
	} else {
		fmt.Println("⚠️  Poor distribution - consider more virtual nodes")
	}

	fmt.Println("\nKey Distribution:")
	for server, count := range metrics.Distribution {
		percentage := float64(count) * 100 / float64(metrics.TotalKeys)
		fmt.Printf("  %s: %d keys (%.1f%%)\n", server, count, percentage)
	}
}

package hashring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetServer(b *testing.B) {
	ring := New(150)
	require.NoError(b, ring.AddServer("server1"))
	require.NoError(b, ring.AddServer("server2"))
	require.NoError(b, ring.AddServer("server3"))
	require.NoError(b, ring.AddServer("server4"))
	require.NoError(b, ring.AddServer("server5"))

	for i := 0; b.Loop(); i++ {
		key := fmt.Sprintf("key-%d", i%10000)
		// NB: We don't do the require step here so it doesn't influence the benchmarks.
		_, _ = ring.GetServer(key)
	}
}

func BenchmarkAddServer(b *testing.B) {
	b.Skip("Not implemented yet")

	for b.Loop() {
		b.StopTimer()
		// ring := New(150)
		b.StartTimer()

		// TODO: Write benchmark test for AddServer.
	}
}

func BenchmarkDistribution(b *testing.B) {
	b.Skip("Not implemented yet")

	ring := New(150)
	require.NoError(b, ring.AddServer("server1"))
	require.NoError(b, ring.AddServer("server2"))
	require.NoError(b, ring.AddServer("server3"))

	// TODO: Generate 1000 keys
	// TODO: Benchmark GetDistribution
}

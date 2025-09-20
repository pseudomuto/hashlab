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
		_, _ = ring.GetServer(key)
	}
}

func BenchmarkAddServer(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		ring := New(150)
		b.StartTimer()

		for j := range 10 {
			_ = ring.AddServer(fmt.Sprintf("server-%d", j))
		}
	}
}

func BenchmarkDistribution(b *testing.B) {
	ring := New(150)
	require.NoError(b, ring.AddServer("server1"))
	require.NoError(b, ring.AddServer("server2"))
	require.NoError(b, ring.AddServer("server3"))

	keys := make([]string, 1000)
	for i := range 1000 {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	for b.Loop() {
		ring.GetDistribution(keys)
	}
}

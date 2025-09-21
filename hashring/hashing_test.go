package hashring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	ring := New(150)
	require.NotNil(t, ring, "New returned nil")

	// More assertions will be added as you implement features
}

func TestAddServer(t *testing.T) {
	t.Skip("AddServer not implemented yet")
}

func TestRemoveServer(t *testing.T) {
	t.Skip("RemoveServer not implemented yet")
}

func TestGetServer(t *testing.T) {
	t.Skip("GetServer not implemented yet")
}

func TestDistribution(t *testing.T) {
	t.Skip("Distribution not implemented yet")
}

func TestConsistency(t *testing.T) {
	t.Skip("Consistency not implemented yet")

	// TODO:
	// * Add 3 servers
	// * Map 1000 keys
	// * Add 4th server
	// * Check how many keys have moved
	// * Assert that the number of moved keys is between 700 and 800 (expected at this point).
}

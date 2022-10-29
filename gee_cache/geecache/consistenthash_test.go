package geecache_test

import (
	"geecache"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	t.Parallel()

	hash := geecache.NewMap(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		require.Equal(t, v, hash.Get(k))
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// Now 27 should map to 8
	require.Equal(t, "8", hash.Get("27"))
}

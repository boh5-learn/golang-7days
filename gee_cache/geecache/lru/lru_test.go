package lru_test

import (
	"geecache/lru"
	"testing"

	"github.com/stretchr/testify/require"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Add_Get(t *testing.T) {
	t.Parallel()

	key1 := "key1"
	key2 := "key2"

	lruCache := lru.New(int64(0), nil)
	lruCache.Add(key1, String("1234"))
	v, ok := lruCache.Get(key1)
	require.True(t, ok)
	require.EqualValues(t, "1234", v)

	v, ok = lruCache.Get(key2)
	require.False(t, ok)
	require.Nil(t, v)
}

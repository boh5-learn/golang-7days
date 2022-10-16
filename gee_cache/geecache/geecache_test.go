package geecache_test

import (
	"fmt"
	"geecache"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetter(t *testing.T) {
	t.Parallel()

	var f geecache.Getter = geecache.GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	key := "key"

	v, err := f.Get(key)
	require.NoError(t, err)
	require.Equal(t, []byte(key), v)
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGroup_Get(t *testing.T) {
	t.Parallel()

	loadCounts := make(map[string]int, len(db))
	gee := geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		// load from geecache.GetterFunc
		view, err := gee.Get(k)
		require.NoError(t, err)
		require.Equal(t, v, view.String())

		// cache hit
		_, err = gee.Get(k)
		require.NoError(t, err)
		require.Equal(t, 1, loadCounts[k])
	}

	// get not exist
	view, err := gee.Get("unknown")
	require.Error(t, err)
	require.Empty(t, view)
}

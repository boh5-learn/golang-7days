package geecache_test

import (
	"fmt"
	"geecache"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	t.Parallel()

	addr := "localhost:9999"
	baseUrl := "http://" + addr + "/_geecache/scores/"

	go func() {
		geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
			func(key string) ([]byte, error) {
				log.Println("[SlowDB] search key", key)
				if v, ok := db[key]; ok {
					return []byte(v), nil
				}
				return nil, fmt.Errorf("%s not exist", key)
			}))

		peers := geecache.NewHTTPPool(addr)
		log.Println("geecache is running at", addr)
		log.Fatal(http.ListenAndServe(addr, peers))
	}()

	key := "Tom"
	resp, err := http.Get(baseUrl + key)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	data := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(data)
	require.EqualError(t, err, "EOF")
	require.Equal(t, db[key], string(data))

	key = "NotExist"
	resp, err = http.Get(baseUrl + key)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

package gee

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePattern(t *testing.T) {
	t.Parallel()

	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	require.True(t, ok)
	ok = reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	require.True(t, ok)
	ok = reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	require.True(t, ok)
}

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestGetRoute(t *testing.T) {
	t.Parallel()

	r := newTestRouter()

	n, ps := r.getRoute("GET", "/hello/boh5")
	require.NotNil(t, n)
	require.Equal(t, "/hello/:name", n.pattern)
	require.Equal(t, "boh5", ps["name"])
}

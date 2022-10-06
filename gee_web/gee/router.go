package gee

import (
	"log"
	"net/http"
	"strings"
)

// router store handlers and implement router functions
type router struct {
	roots    map[string]*node       // store trie root node of GET, POST ... method
	handlers map[string]HandlerFunc // store HandlerFunc of pattern
}

// newRouter returns a new router
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern parse pattern to parts, only one '*' is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// routeKey construct key of router
func routeKey(method, pattern string) string {
	return method + "-" + pattern
}

// addRoute add route to router
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s\n", method, pattern)

	parts := parsePattern(pattern)

	key := routeKey(method, pattern)

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute returns matched node and params
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for idx, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[idx]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[idx:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	if n, params := r.getRoute(c.Method, c.Path); n != nil {
		c.Params = params
		key := routeKey(c.Method, n.pattern)
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

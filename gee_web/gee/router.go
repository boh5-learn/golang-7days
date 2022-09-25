package gee

import (
	"log"
	"net/http"
)

// router store handlers and implement router functions
type router struct {
	handlers map[string]HandlerFunc
}

// newRouter returns a new router
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// routeKey construct key of router
func (r *router) routeKey(method, pattern string) string {
	return method + "-" + pattern
}

// addRoute add route to router
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s\n", method, pattern)
	key := r.routeKey(method, pattern)
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := r.routeKey(c.Method, c.Path)
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

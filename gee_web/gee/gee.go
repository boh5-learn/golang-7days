package gee

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc defined the request handler used by gee
type HandlerFunc func(w http.ResponseWriter, req *http.Request)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router map[string]HandlerFunc
}

// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// routeKey construct key of Engine.router
func (e *Engine) routeKey(method, pattern string) string {
	return method + "-" + pattern
}

// addRoute add route to Engine.router
func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	key := e.routeKey(method, pattern)
	e.router[key] = handler
}

// GET defines the method to add GET request
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// Run start a http server
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := e.routeKey(req.Method, req.URL.Path)
	if handler, ok := e.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

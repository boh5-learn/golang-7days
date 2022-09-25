package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]any

// Context construct data of http.Request and http.ResponseWriter
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

// newContext returns a new Context
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// PostForm returns form data of key
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query returns query string of key
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status set status code of response
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader set header of response
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// String set plain text response
// code: response status code
// format: format string
// values: format values
func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	if _, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...))); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// JSON set json response
// code: response status code
// obj: obj to encode
func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Data set []byte data to response
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	if _, err := c.Writer.Write(data); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// HTML set html response
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if _, err := c.Writer.Write([]byte(html)); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

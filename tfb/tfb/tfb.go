package tfb

import (
	"net/http"
)

// HandlerFunc defines the request handler used by tfb
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
}

// New is the constructor of teb.Engine
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

// GET define the method to add GET request
func (e *Engine) GET(patten string, handler HandlerFunc) {
	e.addRoute("GET", patten, handler)
}

// POST define the method to add POST request
func (e *Engine) POST(patten string, handler HandlerFunc) {
	e.addRoute("POST", patten, handler)
}

// Run define the method to start a http sever
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}

package tfb

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by tfb
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup // engine 作为顶层分组，拥有所有路由分组功能
	router       *router
	groups       []*RouterGroup // store all groups
}

// RouterGroup 路由组结构，用于路由分组
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// New is the constructor of teb.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	// add group to engine groups
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET define the method to add GET request
func (group *RouterGroup) GET(patten string, handler HandlerFunc) {
	group.addRoute("GET", patten, handler)
}

// POST define the method to add POST request
func (group *RouterGroup) POST(patten string, handler HandlerFunc) {
	group.addRoute("POST", patten, handler)
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// Run define the method to start a http sever
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// 接收到一个具体请求时，要判断该请求适用于哪些中间件
// 在这里简单通过 URL 的前缀判断。得到中间件列表后，赋值给 c.handlers
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	e.router.handle(c)
}

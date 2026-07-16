package framework

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler is the signature every controller action uses. Controllers receive
// a framework Request and never touch Gin.
type Handler func(*Request)

// Middleware wraps a request. Call next() to continue the chain; returning
// without calling next() aborts the request (the middleware is expected to
// have written a response).
type Middleware func(req *Request, next func())

// Router registers application routes. It mirrors Laravel's Route facade:
//
//	r.Get("/users", handler)
//	r.Resource("/users", &controllers.UserController{...})
//	admin := r.Group("/admin", middleware.Auth(app))
type Router struct {
	app   *App
	group *gin.RouterGroup
}

func newRouter(app *App, group *gin.RouterGroup) *Router {
	return &Router{app: app, group: group}
}

// Get registers a GET route.
func (r *Router) Get(path string, h Handler) { r.group.GET(path, r.wrap(h)) }

// Post registers a POST route.
func (r *Router) Post(path string, h Handler) { r.group.POST(path, r.wrap(h)) }

// Put registers a PUT route.
func (r *Router) Put(path string, h Handler) { r.group.PUT(path, r.wrap(h)) }

// Patch registers a PATCH route.
func (r *Router) Patch(path string, h Handler) { r.group.PATCH(path, r.wrap(h)) }

// Delete registers a DELETE route.
func (r *Router) Delete(path string, h Handler) { r.group.DELETE(path, r.wrap(h)) }

// Use attaches middleware to this router and everything registered after it.
func (r *Router) Use(mw ...Middleware) {
	for _, m := range mw {
		r.group.Use(r.wrapMiddleware(m))
	}
}

// Group returns a sub-router mounted at prefix with optional middleware.
func (r *Router) Group(prefix string, mw ...Middleware) *Router {
	handlers := make([]gin.HandlerFunc, 0, len(mw))
	for _, m := range mw {
		handlers = append(handlers, r.wrapMiddleware(m))
	}
	return newRouter(r.app, r.group.Group(prefix, handlers...))
}

// resourceActions maps the seven RESTful controller methods to routes,
// exactly like Laravel's Route::resource.
var resourceActions = []struct {
	method string
	path   string
	action string
}{
	{"GET", "", "Index"},
	{"GET", "/create", "Create"},
	{"POST", "", "Store"},
	{"GET", "/:id", "Show"},
	{"GET", "/:id/edit", "Edit"},
	{"PUT", "/:id", "Update"},
	{"PATCH", "/:id", "Update"},
	{"DELETE", "/:id", "Destroy"},
}

// Resource registers RESTful routes for a controller struct. Any of the
// seven actions (Index, Create, Store, Show, Edit, Update, Destroy) the
// controller implements — with signature func(*framework.Request) — is
// registered; missing actions are simply skipped.
func (r *Router) Resource(path string, controller any) {
	value := reflect.ValueOf(controller)
	registered := 0

	for _, def := range resourceActions {
		method := value.MethodByName(def.action)
		if !method.IsValid() {
			continue
		}
		handler, ok := method.Interface().(func(*Request))
		if !ok {
			panic(fmt.Sprintf(
				"framework: %s.%s must have signature func(*framework.Request)",
				value.Type(), def.action,
			))
		}
		r.handleMethod(def.method, path+def.path, handler)
		registered++
	}

	if registered == 0 {
		panic(fmt.Sprintf(
			"framework: Resource(%q, %s): controller implements none of the resource actions",
			path, value.Type(),
		))
	}
}

func (r *Router) handleMethod(method, path string, h Handler) {
	switch strings.ToUpper(method) {
	case "GET":
		r.Get(path, h)
	case "POST":
		r.Post(path, h)
	case "PUT":
		r.Put(path, h)
	case "PATCH":
		r.Patch(path, h)
	case "DELETE":
		r.Delete(path, h)
	}
}

const requestContextKey = "flux.request"

// requestFrom returns the per-request *Request, creating it on first use so
// middleware and the final handler share the same instance (and therefore
// the same authenticated user, parsed body, etc.).
func requestFrom(app *App, c *gin.Context) *Request {
	if existing, ok := c.Get(requestContextKey); ok {
		return existing.(*Request)
	}
	req := newRequest(app, c)
	c.Set(requestContextKey, req)
	return req
}

func (r *Router) wrap(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(requestFrom(r.app, c))
	}
}

func (r *Router) wrapMiddleware(mw Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := requestFrom(r.app, c)
		nextCalled := false
		mw(req, func() {
			nextCalled = true
			c.Next()
		})
		if !nextCalled {
			c.Abort()
		}
	}
}

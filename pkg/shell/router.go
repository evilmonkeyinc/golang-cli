package shell

import (
	"strings"

	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
)

// The Router interface details the core shell router functions.
type Router interface {
	Handler
	Routes
	// Use appends one or more middleware onto the router stack.
	Use(...Middleware)
	// Group adds a new inline-router to the router stack.
	Group(func(r Router)) Router
	// Route adds a new sub-router to the router stack, along the specified command path.
	Route(string, func(r Router)) Router
	// Handle adds a shell handler to the router stack, along the specified command path.
	Handle(string, Handler)
	// HandleFunction adds a shell handler function to the router stack, along the specified command path.
	HandleFunction(string, HandlerFunction)
	// NotFound defines a shell handler that will respond if a command path cannot be evaluated.
	NotFound(Handler)
}

// Routes interface describes functions for router traversal.
type Routes interface {
	// Routes returns the linked shell handlers.
	Routes() map[string]Handler
	// Middlewares returns the list of middlewares in use by the router.
	Middlewares() []Middleware
	// Match evaluates the routing tree for a handler that matches the supplied arguments
	// and returns the handler, wrapped in the appropriate middleware handler functions
	Match([]string) (Handler, bool)
}

// newRouter will return a new empty router
func newRouter() *router {
	return &router{
		children:        []Router{},
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		parent:          nil,
		notFoundHandler: nil,
	}
}

// childRouter will create a new sub router as an inline group router
func childRouter(rtr *router) *router {
	return &router{
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		parent:          rtr,
		children:        []Router{},
		notFoundHandler: rtr.notFoundHandler,
	}
}

// subRouter will create a new sub router as a sub command router
func subRouter(rtr *router) *router {
	return &router{
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		parent:          nil,
		children:        []Router{},
		notFoundHandler: rtr.notFoundHandler,
	}
}

type router struct {
	handlers   map[string]Handler
	middleware []Middleware

	parent   Router
	children []Router

	notFoundHandler Handler
}

func (rtr *router) Execute(writer ResponseWriter, request *Request) error {
	args := request.Args
	if handler, found := rtr.Match(args); found {
		request = request.WithRoutes(args[0], rtr)
		return handler.Execute(writer, request)
	}

	if rtr.notFoundHandler != nil {
		handler := chain(rtr.middleware, rtr.notFoundHandler)
		return handler.Execute(writer, request)

	}

	return nil
}

func (rtr *router) Routes() map[string]Handler {
	return rtr.handlers
}

func (rtr *router) Middlewares() []Middleware {
	return rtr.middleware
}

func (rtr *router) Match(args []string) (Handler, bool) {
	if len(args) == 0 {
		return nil, false
	}

	arg := args[0]
	for key, handler := range rtr.handlers {
		if strings.EqualFold(arg, key) {
			return chain(rtr.middleware, handler), true
		}
	}

	for _, sub := range rtr.children {
		if handler, found := sub.Match(args); found {
			return chain(rtr.middleware, handler), true
		}
	}

	return nil, false
}

func (rtr *router) Use(middleware ...Middleware) {
	rtr.middleware = append(rtr.middleware, middleware...)
}

func (rtr *router) Group(setup func(r Router)) Router {
	subRouter := childRouter(rtr)
	setup(subRouter)
	rtr.children = append(rtr.children, subRouter)
	return subRouter
}

func (rtr *router) Route(command string, setup func(r Router)) Router {
	if _, exists := rtr.Match([]string{command}); exists {
		panic(errors.DuplicateCommand(command))
	}

	subRouter := subRouter(rtr)
	setup(subRouter)
	rtr.handlers[command] = subRouter
	return subRouter
}

func (rtr *router) Handle(command string, handler Handler) {
	if _, exists := rtr.Match([]string{command}); exists {
		panic(errors.DuplicateCommand(command))
	}
	rtr.handlers[command] = handler
}

func (rtr *router) HandleFunction(command string, handerFunction HandlerFunction) {
	if _, exists := rtr.Match([]string{command}); exists {
		panic(errors.DuplicateCommand(command))
	}
	rtr.handlers[command] = handerFunction
}

func (rtr *router) NotFound(handler Handler) {
	rtr.notFoundHandler = handler
}

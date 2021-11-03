package shell

import (
	"strings"

	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
)

type Router interface {
	Handler
	Routes
	Use(...Middleware)
	Group(func(r Router)) Router
	Route(string, func(r Router)) Router
	Handle(string, Handler)
	HandleFunction(string, HandlerFunction)
	Help(Handler)
	NotFound(Handler)
}

type Routes interface {
	Routes() map[string]Handler
	Middlewares() []Middleware
	Match([]string) (Handler, bool)
}

// newRouter will return a new empty router
func newRouter() *router {
	return &router{
		children:        []Router{},
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		parent:          nil,
		helpHander:      nil,
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
		helpHander:      rtr.helpHander,
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
		helpHander:      rtr.helpHander,
		notFoundHandler: rtr.notFoundHandler,
	}
}

type router struct {
	handlers   map[string]Handler
	middleware []Middleware

	parent   Router
	children []Router

	helpHander      Handler
	notFoundHandler Handler
}

func (rtr *router) Execute(writer ResponseWriter, request Request) error {
	// TODO : flag logic
	// needs to use interface so can support pflags AND standard flags

	args := request.Args()
	if handler, found := rtr.Match(args); found {
		request = request.WithArgs(args[1:]).WithRoutes(rtr)
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

	// TODO : check args if this is a help request
	// if help arg then return current help
	// if help flag, identify possible sub router and call help there
	// if help is called but a handler is also found, then we need to
	// let help have access to it so it can display command specific help information
	// Should help be middleware?

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

func (rtr *router) Help(handler Handler) {
	rtr.helpHander = handler
}

func (rtr *router) NotFound(handler Handler) {
	rtr.notFoundHandler = handler
}

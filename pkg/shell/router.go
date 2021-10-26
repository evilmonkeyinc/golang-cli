package shell

import (
	"fmt"
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
	Help(HandlerFunction)
	NotFound(HandlerFunction)
}

type Routes interface {
	Routes() map[string]Handler
	Middlewares() []Middleware
	Match([]string) (Handler, bool)
}

// newRouter will return a new empty router
func newRouter() *router {
	return &router{
		subRouters: []Router{},
		handlers:   map[string]Handler{},
		middleware: []Middleware{},
	}
}

// subRouter will create a new sub router for the supplied router
func subRouter(rtr *router) *router {
	return &router{
		subRouters:      []Router{},
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		helpHander:      rtr.helpHander,
		notFoundHandler: rtr.notFoundHandler,
	}
}

type router struct {
	subRouters []Router
	handlers   map[string]Handler
	middleware []Middleware

	helpHander      HandlerFunction
	notFoundHandler HandlerFunction
}

func (rtr *router) Execute(writer ResponseWriter, request Request) error {
	// TODO : flag logic

	args := request.Args()
	routes := request.Routes()

	var toExecute Handler
	if handler, found := rtr.Match(request.Args()); found {
		args = request.Args()[1:]
		toExecute = handler
		routes = rtr
	}

	if toExecute == nil && rtr.notFoundHandler != nil {
		toExecute = rtr.notFoundHandler
	}
	if toExecute == nil && rtr.helpHander != nil {
		toExecute = rtr.helpHander
	}

	if toExecute != nil {
		request = request.WithArgs(args).WithRoutes(routes)

		if len(rtr.middleware) > 0 {
			handlerFunc := toExecute
			handlerFunc = rtr.middleware[len(rtr.middleware)-1](handlerFunc)
			for i := len(rtr.middleware) - 2; i >= 0; i-- {
				handlerFunc = rtr.middleware[i](handlerFunc)
			}
			return handlerFunc.Execute(writer, request)
		} else {
			return toExecute.Execute(writer, request)
		}
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
			return handler, true
		}
	}

	for _, sub := range rtr.subRouters {
		if handler, found := sub.Match(args); found {
			return handler, true
		}
	}

	return nil, false
}

func (rtr *router) Use(middleware ...Middleware) {
	rtr.middleware = append(rtr.middleware, middleware...)
}

func (rtr *router) Group(setup func(r Router)) Router {
	subRouter := subRouter(rtr)
	setup(subRouter)
	rtr.subRouters = append(rtr.subRouters, subRouter)
	return subRouter
}

func (rtr *router) Route(command string, setup func(r Router)) Router {
	if _, exists := rtr.handlers[command]; exists {
		panic(fmt.Errorf("%s %w", command, errors.DuplicateCommandFunction))
	}

	subRouter := subRouter(rtr)
	setup(subRouter)
	rtr.handlers[command] = subRouter
	return subRouter
}

func (rtr *router) Handle(command string, handler Handler) {
	if _, exists := rtr.handlers[command]; exists {
		panic(fmt.Errorf("%s %w", command, errors.DuplicateCommandFunction))
	}
	rtr.handlers[command] = handler
}

func (rtr *router) HandleFunction(command string, handerFunction HandlerFunction) {
	if _, exists := rtr.handlers[command]; exists {
		panic(fmt.Errorf("%s %w", command, errors.DuplicateCommandFunction))
	}
	rtr.handlers[command] = handerFunction
}

func (rtr *router) Help(handler HandlerFunction) {
	rtr.helpHander = handler
}

func (rtr *router) NotFound(handler HandlerFunction) {
	rtr.notFoundHandler = handler
}

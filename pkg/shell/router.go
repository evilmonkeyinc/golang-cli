package shell

import (
	"fmt"
	"strings"

	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
)

// The Router interface details the core shell router functions.
type Router interface {
	Handler
	Routes
	// Flags adds a FlagHandler that will add flags to the request FlagSet before
	// it attempts to match a command.
	Flags(FlagHandler)
	// Group adds a new inline-router to the router stack.
	Group(func(r Router)) Router
	// Handle adds a shell handler to the router stack, along the specified command path.
	Handle(string, Handler)
	// HandleFunction adds a shell handler function to the router stack, along the specified command path.
	HandleFunction(string, HandlerFunction)
	// Route adds a new sub-router to the router stack, along the specified command path.
	NotFound(Handler)
	// Use appends one or more middleware onto the router stack.
	Route(string, func(r Router)) Router
	// NotFound defines a shell handler that will respond if a command path cannot be evaluated.
	Use(...Middleware)
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
		notFoundHandler: nil,
		parent:          nil,
	}
}

// childRouter will create a new sub router as an inline group router
func childRouter(rtr *router) *router {
	return &router{
		children:        []Router{},
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		notFoundHandler: rtr.notFoundHandler,
		parent:          rtr,
	}
}

// subRouter will create a new sub router as a sub command router
func subRouter(rtr *router) *router {
	return &router{
		children:        []Router{},
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		notFoundHandler: rtr.notFoundHandler,
		parent:          nil,
	}
}

type router struct {
	children        []Router
	flags           FlagHandler
	handlers        map[string]Handler
	middleware      []Middleware
	notFoundHandler Handler
	parent          Router
}

func (rtr *router) Execute(writer ResponseWriter, request *Request) error {
	args := request.Args
	flagSet := request.FlagSet

	if handler, found := rtr.Match(args); found {
		currentRoute := args[0]
		flagSet = flagSet.SubFlagSet(currentRoute)
		if flagHandler, ok := handler.(FlagHandler); ok {
			flagHandler.Define(flagSet)
		}
		var parseErr error = nil
		if args, parseErr = flagSet.Parse(args[1:]); parseErr != nil {
			// TODO : check for ErrHelp
			fmt.Fprintln(writer.ErrorWriter(), parseErr.Error())
		}
		request = request.UpdateRequest(currentRoute, args, flagSet, rtr)
		return handler.Execute(writer, request)
	}

	if rtr.notFoundHandler != nil {
		handler := &chainHandler{
			handler:     rtr.notFoundHandler,
			middlewares: rtr.middleware,
		}
		request = request.UpdateRequest("", args, flagSet, rtr)
		return handler.Execute(writer, request)
	}

	return nil
}

func (rtr *router) Flags(fn FlagHandler) {
	rtr.flags = fn
}

func (rtr *router) Define(fd FlagDefiner) {
	if rtr.flags != nil {
		rtr.flags.Define(fd)
	}
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
			return &chainHandler{
				handler:     handler,
				middlewares: rtr.middleware,
			}, true
		}
	}

	for _, sub := range rtr.children {
		if handler, found := sub.Match(args); found {
			return &chainHandler{
				handler:     handler,
				middlewares: rtr.middleware,
			}, true
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

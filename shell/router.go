package shell

import (
	"fmt"
	"strings"

	"github.com/evilmonkeyinc/golang-cli/errors"
	"github.com/evilmonkeyinc/golang-cli/flags"
)

// The Router interface details the core shell router functions.
type Router interface {
	Handler
	Routes
	// Flags adds a FlagHandler that will add flags to the request FlagSet before
	// it attempts to match a command.
	Flags(flags.FlagHandler)
	// Group adds a new inline-router to the router stack.
	Group(func(r Router)) Router
	// Handle adds a shell handler to the router stack, along the specified command path.
	Handle(string, Handler)
	// HandleFunction adds a shell handler function to the router stack, along the specified command path.
	HandleFunction(string, HandlerFunction)
	// NotFound defines a shell handler that will respond if a command path cannot be evaluated.
	NotFound(Handler)
	// Route adds a new sub-router to the router stack, along the specified command path.
	Route(string, func(r Router)) Router
	// Mount adds the specified router to the router stack, along the specified command path.
	//
	// A mounted router will not inherit helper functions, such as the not found handler,
	// from the parent router in the same way a sub-router created by the Route() does, you must
	// set these manually.
	Mount(string, Router)
	// Use appends one or more middleware onto the router stack.
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
func newRouter() *StandardRouter {
	return &StandardRouter{
		children:        []Router{},
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		notFoundHandler: nil,
		parent:          nil,
	}
}

// childRouter will create a new sub router as an inline group router
func childRouter(rtr *StandardRouter) *StandardRouter {
	return &StandardRouter{
		children:        []Router{},
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		notFoundHandler: rtr.notFoundHandler,
		parent:          rtr,
	}
}

// subRouter will create a new sub router as a sub command router
func subRouter(rtr *StandardRouter) *StandardRouter {
	return &StandardRouter{
		children:        []Router{},
		handlers:        map[string]Handler{},
		middleware:      []Middleware{},
		notFoundHandler: rtr.notFoundHandler,
		parent:          nil,
	}
}

// StandardRouter represents the standard implementation of the Router interface.
type StandardRouter struct {
	children        []Router
	flags           flags.FlagHandler
	handlers        map[string]Handler
	middleware      []Middleware
	notFoundHandler Handler
	parent          Router
}

func (rtr *StandardRouter) setup() {
	if rtr.handlers == nil {
		rtr.handlers = make(map[string]Handler)
	}
}

// Execute is used to execute the shell handler.
func (rtr *StandardRouter) Execute(writer ResponseWriter, request *Request) error {
	args := request.Args
	flagSet := request.FlagSet

	if handler, found := rtr.Match(args); found {
		currentRoute := args[0]
		flagSet = flagSet.SubFlagSet(currentRoute)
		if flagHandler, ok := handler.(flags.FlagHandler); ok {
			flagHandler.Define(flagSet)
		}
		var parseErr error = nil
		if args, parseErr = flagSet.Parse(args[1:]); parseErr != nil {
			if errors.IsHelpRequested(parseErr) {
				return parseErr
			}
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

// Define allows the function to define command-line flags.
func (rtr *StandardRouter) Define(fd flags.FlagDefiner) {
	if rtr.flags != nil {
		rtr.flags.Define(fd)
	}
}

// Routes returns the linked shell handlers.
func (rtr *StandardRouter) Routes() map[string]Handler {
	return rtr.handlers
}

// Middlewares returns the list of middlewares in use by the router.
func (rtr *StandardRouter) Middlewares() []Middleware {
	return rtr.middleware
}

// Match evaluates the routing tree for a handler that matches the supplied arguments
// and returns the handler, wrapped in the appropriate middleware handler functions
func (rtr *StandardRouter) Match(args []string) (Handler, bool) {
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

// Flags adds a FlagHandler that will add flags to the request FlagSet before
// it attempts to match a command.
func (rtr *StandardRouter) Flags(fn flags.FlagHandler) {
	rtr.flags = fn
}

// Group adds a new inline-router to the router stack.
func (rtr *StandardRouter) Group(setup func(r Router)) Router {
	subRouter := childRouter(rtr)
	setup(subRouter)
	rtr.children = append(rtr.children, subRouter)
	return subRouter
}

// Handle adds a shell handler to the router stack, along the specified command path.
func (rtr *StandardRouter) Handle(command string, handler Handler) {
	rtr.setup()
	if _, exists := rtr.Match([]string{command}); exists {
		panic(errors.DuplicateCommand(command))
	}
	rtr.handlers[command] = handler
}

// HandleFunction adds a shell handler function to the router stack, along the specified command path.
func (rtr *StandardRouter) HandleFunction(command string, handerFunction HandlerFunction) {
	rtr.setup()
	if _, exists := rtr.Match([]string{command}); exists {
		panic(errors.DuplicateCommand(command))
	}
	rtr.handlers[command] = handerFunction
}

// NotFound defines a shell handler that will respond if a command path cannot be evaluated.
func (rtr *StandardRouter) NotFound(handler Handler) {
	rtr.notFoundHandler = handler
}

// Route adds a new sub-router to the router stack, along the specified command path.
func (rtr *StandardRouter) Route(command string, setup func(r Router)) Router {
	rtr.setup()
	if _, exists := rtr.Match([]string{command}); exists {
		panic(errors.DuplicateCommand(command))
	}

	subRouter := subRouter(rtr)
	setup(subRouter)
	rtr.handlers[command] = subRouter
	return subRouter
}

// Mount adds the specified router to the router stack, along the specified command path.
//
// A mounted router will not inherit helper functions, such as the not found handler,
// from the parent router in the same way a sub-router created by the Route() does, you must
// set these manually.
func (rtr *StandardRouter) Mount(command string, router Router) {
	rtr.setup()
	if _, exists := rtr.Match([]string{command}); exists {
		panic(errors.DuplicateCommand(command))
	}

	rtr.handlers[command] = router
}

// Use appends one or more middleware onto the router stack.
func (rtr *StandardRouter) Use(middleware ...Middleware) {
	rtr.middleware = append(rtr.middleware, middleware...)
}

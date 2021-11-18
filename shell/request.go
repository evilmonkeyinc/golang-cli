package shell

import (
	"context"
	"strings"
)

// NewRequest wraps NewRequestWithContext using context.Background.
func NewRequest(path, args []string, flagSet FlagSet, routes Routes) *Request {
	return NewRequestWithContext(context.Background(), path, args, flagSet, routes)
}

// NewRequestWithContext returns a new Request given a path, args, and routes.
func NewRequestWithContext(ctx context.Context, path, args []string, flagSet FlagSet, routes Routes) *Request {
	return &Request{
		ctx:     ctx,
		Args:    args,
		FlagSet: flagSet,
		Path:    path,
		Routes:  routes,
	}
}

// A Request represents a the request sent by the shell and processed by the router and handlers.
type Request struct {
	ctx context.Context

	// Args contains the arguments passed as part of the request.
	Args []string
	// Flagset contains the flagset used to parse arguments.
	FlagSet FlagSet
	// Path contains the request path.
	Path []string
	// Routes contains the router routes functions linked to the executed router.
	Routes Routes
}

// Context returns the request's context. To change the context, use WithContext.
func (request *Request) Context() context.Context {
	return request.ctx
}

// FlagValues returns the parsed flag values for the request flagset.
func (request *Request) FlagValues() FlagValues {
	return request.FlagSet
}

// WithContext returns a shallow copy of the request with its context changed to ctx.
func (request *Request) WithContext(ctx context.Context) *Request {
	args := make([]string, len(request.Args))
	copy(args, request.Args)

	path := make([]string, len(request.Path))
	copy(path, request.Path)

	return &Request{
		ctx:     ctx,
		Args:    args,
		FlagSet: request.FlagSet,
		Path:    path,
		Routes:  request.Routes,
	}
}

// UpdateRequest returns a shallow copy of the request with updated path, args, flagset, and routes.
func (request *Request) UpdateRequest(selectedRoute string, args []string, flagSet FlagSet, routes Routes) *Request {
	if args == nil {
		args = make([]string, len(request.Args))
		copy(args, request.Args)
	}

	path := make([]string, len(request.Path))
	copy(path, request.Path)

	if selectedRoute != "" {
		path = append(path, selectedRoute)
		if len(args) > 0 && strings.EqualFold(args[0], selectedRoute) {
			args = args[1:]
		}
	}

	if routes == nil {
		routes = request.Routes
	}

	if flagSet == nil {
		flagSet = request.FlagSet
	}

	return &Request{
		ctx:     request.ctx,
		Args:    args,
		FlagSet: flagSet,
		Path:    path,
		Routes:  routes,
	}
}

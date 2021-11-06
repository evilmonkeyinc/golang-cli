package shell

import (
	"context"
	"strings"
)

// NewRequest wraps NewRequestWithContext using context.Background.
func NewRequest(path, args []string, routes Routes) *Request {
	return NewRequestWithContext(context.Background(), path, args, routes)
}

// NewRequestWithContext returns a new Request given a path, args, and routes.
func NewRequestWithContext(ctx context.Context, path, args []string, routes Routes) *Request {
	return &Request{
		ctx:    ctx,
		Args:   args,
		Path:   path,
		Routes: routes,
	}
}

// A Request represents a the request sent by the shell and processed by the router and handlers.
type Request struct {
	ctx context.Context

	// Args contains the arguments passed as part of the request.
	Args []string
	// Path contains the request path.
	Path []string
	// Routes contains the router routes functions linked to the executed router.
	Routes Routes
}

// Context returns the request's context. To change the context, use WithContext.
func (request *Request) Context() context.Context {
	return request.ctx
}

// WithContext returns a shallow copy of the request with its context changed to ctx.
func (request *Request) WithContext(ctx context.Context) *Request {
	args := make([]string, len(request.Args))
	copy(args, request.Args)

	path := make([]string, len(request.Path))
	copy(path, request.Path)

	return &Request{
		ctx:    ctx,
		Args:   args,
		Path:   path,
		Routes: request.Routes,
	}
}

// WithRoutes returns a shallow copy of the request with updated Routes,
// the selected route added to paths array, and removed from the args array
func (request *Request) WithRoutes(selectedRoute string, routes Routes) *Request {
	args := make([]string, len(request.Args))
	copy(args, request.Args)

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

	return &Request{
		ctx:    request.ctx,
		Args:   args,
		Path:   path,
		Routes: routes,
	}
}

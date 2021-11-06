package shell

// The MiddlewareFunction type is an adapter to allow the use of ordinary functions as shell middleware.
type MiddlewareFunction func(next Handler) Handler

// Handle is used to perform the middleware function.
func (middleware MiddlewareFunction) Handle(next Handler) Handler {
	return middleware(next)
}

// The Middleware interface describes a shell middleware function
type Middleware interface {
	// Handle is used to perform the middleware function.
	Handle(next Handler) Handler
}

func chain(middlewares []Middleware, handler Handler) Handler {
	if len(middlewares) == 0 {
		return handler
	}

	chainedHandler := handler
	chainedHandler = middlewares[len(middlewares)-1].Handle(chainedHandler)
	for i := len(middlewares) - 2; i >= 0; i-- {
		chainedHandler = middlewares[i].Handle(chainedHandler)
	}
	return chainedHandler
}

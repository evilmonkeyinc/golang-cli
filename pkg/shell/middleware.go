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

type chainHandler struct {
	handler     Handler
	middlewares []Middleware
}

func (chain *chainHandler) Define(flagDefiner FlagDefiner) {
	if flagHandler, ok := chain.handler.(FlagHandler); ok {
		flagHandler.Define(flagDefiner)
	}
}

func (chain *chainHandler) Execute(rw ResponseWriter, r *Request) error {
	return chain.chain().Execute(rw, r)
}

func (chain *chainHandler) chain() Handler {
	chainMiddleware := chain.middlewares
	if len(chainMiddleware) == 0 {
		return chain.handler
	}

	chainedHandler := chain.handler
	chainedHandler = chainMiddleware[len(chainMiddleware)-1].Handle(chainedHandler)
	for i := len(chainMiddleware) - 2; i >= 0; i-- {
		chainedHandler = chainMiddleware[i].Handle(chainedHandler)
	}
	return chainedHandler
}

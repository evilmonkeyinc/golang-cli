package shell

type MiddlewareFunction func(next Handler) Handler

func (middleware MiddlewareFunction) Handle(next Handler) Handler {
	return middleware(next)
}

type Middleware interface {
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

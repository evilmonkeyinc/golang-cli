package shell

// Validate the MiddlewareFunction func matches the Middleware interface
var _ Middleware = MiddlewareFunction(func(next Handler) Handler {
	return next
})

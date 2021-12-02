package middleware

import (
	"context"

	"github.com/evilmonkeyinc/golang-cli/shell"
)

const (
	ctxValueError ctxValue = "error"
)

// IsError is used by the ErrorHandler middleware to determine if the
// received error should be used to execute the the ErrorHandlerFunction.
type IsError func(error) bool

// ErrorHandler returns middleware that will intercept all errors returned
// and execute a handler if the error matches a specified condition.
func ErrorHandler(validate IsError, handler shell.Handler) shell.Middleware {
	return shell.MiddlewareFunction(func(next shell.Handler) shell.Handler {
		return shell.HandlerFunction(func(rw shell.ResponseWriter, r *shell.Request) error {
			if err := next.Execute(rw, r); err != nil {
				if validate(err) {
					ctx := r.Context()
					ctx = context.WithValue(ctx, ctxValueError, err)
					r = r.WithContext(ctx)
					return handler.Execute(rw, r)
				}
				return err
			}

			return nil
		})
	})
}

// ErrorHandlerCaughtError returns the error caught by the
// ErrorHandler middleware to be referenced during the middleware handler.
func ErrorHandlerCaughtError(ctx context.Context) error {
	if err, ok := ctx.Value(ctxValueError).(error); ok {
		return err
	}
	return nil
}

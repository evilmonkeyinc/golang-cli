package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
)

// Recoverer returns middleware which will allow the shell to recover from a panic
// The middleware will output the error and stack trace to the ResponseWriter error writer.
func Recoverer() shell.Middleware {
	return shell.MiddlewareFunction(func(next shell.Handler) shell.Handler {
		return shell.HandlerFunction(func(rw shell.ResponseWriter, r *shell.Request) error {
			defer func() {
				if rvr := recover(); rvr != nil {
					errMsg := fmt.Sprintf("%v", rvr)
					switch v := rvr.(type) {
					case error:
						errMsg = v.Error()
					case string:
						errMsg = v
					default:
						break
					}
					fmt.Fprintln(rw.ErrorWriter(), errMsg)
					fmt.Fprintln(rw.ErrorWriter(), string(debug.Stack()))
				}
			}()

			return next.Execute(rw, r)
		})
	})
}

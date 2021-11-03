package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
)

func Recoverer() shell.Middleware {
	return shell.MiddlewareFunction(func(next shell.Handler) shell.Handler {
		return shell.HandlerFunction(func(rw shell.ResponseWriter, r shell.Request) error {
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
					rw.WriteError([]byte(errMsg))
					rw.WriteError([]byte("\n"))
					rw.WriteError(debug.Stack())
				}
			}()

			return next.Execute(rw, r)
		})
	})
}

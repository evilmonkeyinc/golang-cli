package middleware

import (
	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
)

func HandlerFunction(fn shell.HandlerFunction) shell.Handler {
	return &handlerWrapper{
		handlerFunc: fn,
	}
}

type handlerWrapper struct {
	handlerFunc shell.HandlerFunction
}

func (wrapper *handlerWrapper) Execute(rw shell.ResponseWriter, r shell.Request) error {
	return wrapper.handlerFunc(rw, r)
}

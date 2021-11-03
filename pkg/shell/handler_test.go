package shell

import "fmt"

// Validate the HandlerFunction func matches the Handler interface
var _ Handler = HandlerFunction(func(rw ResponseWriter, r Request) error {
	return nil
})

type testHandler struct {
	message string
}

func (handler *testHandler) Execute(ResponseWriter, Request) error {
	return fmt.Errorf(handler.message)
}

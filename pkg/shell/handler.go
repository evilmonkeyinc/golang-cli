package shell

// The HandlerFunction type is an adapter to allow the use of ordinary functions as shell handlers.
type HandlerFunction func(ResponseWriter, *Request) error

// Execute is used to execute the shell handler function.
func (fn HandlerFunction) Execute(writer ResponseWriter, request *Request) error {
	return fn(writer, request)
}

// The Handler interface describes a shell handler functions.
type Handler interface {
	// Execute is used to execute the shell handler.
	Execute(ResponseWriter, *Request) error
}

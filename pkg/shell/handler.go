package shell

type HandlerFunction func(ResponseWriter, Request) error

func (fn HandlerFunction) Execute(writer ResponseWriter, request Request) error {
	return fn(writer, request)
}

type Handler interface {
	Execute(ResponseWriter, Request) error
}

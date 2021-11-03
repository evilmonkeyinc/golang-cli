package shell

import (
	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
)

type Command struct {
	Name        string
	Summary     string
	Description string
	Function    HandlerFunction
}

func (command *Command) Execute(writer ResponseWriter, request Request) error {
	if command.Function == nil {
		return errors.CommandNotFound(command.Name)
	}

	return command.Function(writer, request)
}

package commands

import (
	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
)

// Command allows to define a shell handler with metadata.
type Command struct {
	Name        string
	Summary     string
	Description string
	Function    shell.HandlerFunction
}

// Execute will execute the command function
func (command *Command) Execute(writer shell.ResponseWriter, request *shell.Request) error {
	if command.Function == nil {
		return errors.CommandNotFound(command.Name)
	}

	return command.Function(writer, request)
}

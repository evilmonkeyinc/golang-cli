package commands

import (
	"github.com/evilmonkeyinc/golang-cli/errors"
	"github.com/evilmonkeyinc/golang-cli/shell"
)

// Command allows to define a shell handler with metadata.
type Command struct {
	Name        string
	Summary     string
	Description string
	Flags       shell.FlagHandlerFunction
	Function    shell.HandlerFunction
}

// Define allows the function to define command-line
func (command *Command) Define(flagDefiner shell.FlagDefiner) {
	if command.Flags != nil {
		command.Flags(flagDefiner)
	}
}

// Execute will execute the command function
func (command *Command) Execute(writer shell.ResponseWriter, request *shell.Request) error {
	if command.Function == nil {
		return errors.CommandNotFound(command.Name)
	}

	return command.Function(writer, request)
}

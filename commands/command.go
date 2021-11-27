package commands

import (
	"github.com/evilmonkeyinc/golang-cli/errors"
	"github.com/evilmonkeyinc/golang-cli/flags"
	"github.com/evilmonkeyinc/golang-cli/shell"
)

// Command allows to define a shell handler compatible with the HelpCommand.
type Command struct {
	// The name of the command.
	Name string
	// A short summary of the command.
	Summary string
	// A long description of the command.
	Description string
	// An example of the command used to execute the command.
	Usage string
	// An optional function to include flag definition to the command.
	Flags flags.FlagHandlerFunction
	// The shell handler function to be executed for the command.
	Function shell.HandlerFunction
}

// GetName returns the name of the command handler.
func (command *Command) GetName() string {
	return command.Name
}

// GetSummary returns the short summary of the command handler.
func (command *Command) GetSummary() string {
	return command.Summary
}

// GetDescription returns the long description of the command handler.
func (command *Command) GetDescription() string {
	return command.Description
}

// GetUsage returns an example of the command used to execute the command.
func (command *Command) GetUsage() string {
	return command.Usage
}

// Define allows the function to define command-line
func (command *Command) Define(flagDefiner flags.FlagDefiner) {
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

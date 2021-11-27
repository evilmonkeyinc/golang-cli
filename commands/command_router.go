package commands

import (
	"github.com/evilmonkeyinc/golang-cli/shell"
)

// NewCommandRouter creates a new CommandRouter and runs the setup function for the new router.
func NewCommandRouter(name, summary, description, usage string, setup func(shell.Router)) *CommandRouter {
	router := &CommandRouter{
		Name:        name,
		Summary:     summary,
		Description: description,
		Usage:       usage,
	}
	if setup != nil {
		setup(router)
	}

	return router
}

// CommandRouter is a StandardRouter that conforms to the CommandHandler interface.
// This allows you to define a group of sub commands which are compatible with the
// HelpCommand function.
type CommandRouter struct {
	shell.StandardRouter

	// The name of the command.
	Name string
	// A short summary of the command.
	Summary string
	// A long description of the command.
	Description string
	// An example of the command used to execute the command.
	Usage string
}

// GetName returns the name of the command handler.
func (router *CommandRouter) GetName() string {
	return router.Name
}

// GetSummary returns the short summary of the command handler.
func (router *CommandRouter) GetSummary() string {
	return router.Summary
}

// GetDescription returns the long description of the command handler.
func (router *CommandRouter) GetDescription() string {
	return router.Description
}

// GetUsage returns an example of the command used to execute the command.
func (command *CommandRouter) GetUsage() string {
	return command.Usage
}

package commands

import (
	"github.com/evilmonkeyinc/golang-cli/flags"
	"github.com/evilmonkeyinc/golang-cli/shell"
)

// CommandHandler describes a shell function handler that
// can be described by the HelpCommand
type CommandHandler interface {
	flags.FlagHandler
	shell.Handler
	// GetName returns the name of the command handler.
	GetName() string
	// GetSummary returns the short summary of the command handler.
	GetSummary() string
	// GetDescription returns the long description of the command handler.
	GetDescription() string
	// GetUsage returns an example of the command used to execute the command.
	GetUsage() string
}

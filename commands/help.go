package commands

import (
	"fmt"

	"github.com/evilmonkeyinc/golang-cli/shell"
)

// HelpCommand allows for basic output of command metadata.
type HelpCommand struct {
	// The command string that should be used to call this function,
	// used in the default function output
	Usage string
}

// Execute will execute the help command
func (command *HelpCommand) Execute(writer shell.ResponseWriter, request *shell.Request) error {
	routes := request.Routes

	if command.Usage == "" && len(request.Path) > 0 {
		command.Usage = request.Path[len(request.Path)-1]
	}

	commands := make(map[string]*Command)
	for cmdName, handler := range routes.Routes() {
		if cmd, ok := handler.(*Command); ok {
			commands[cmdName] = cmd
		}
	}

	args := request.Args
	if len(args) > 0 {
		cmdName := args[0]
		if cmd, ok := commands[cmdName]; ok {
			fmt.Fprintf(writer, "\n%s\n", cmdName)
			fmt.Fprintf(writer, "  Usage: %s\n", cmd.Name)
			fmt.Fprintf(writer, "  %s\n\n", cmd.Summary)
			fmt.Fprintf(writer, "%s\n\n", cmd.Description)
			return nil
		}
	}

	fmt.Fprintf(writer, "\n%s: %s\n", command.Usage, fmt.Sprintf("%s or %s <command-name>", command.Usage, command.Usage))
	fmt.Fprintln(writer, "\nAvailable commands")
	fmt.Fprintln(writer, "------------------")
	for cmdName, cmd := range commands {
		fmt.Fprintf(writer, "%12s:\t%s\n", cmdName, cmd.Summary)
	}
	fmt.Fprintf(writer, "\nUse \"%s <command-name>\" for detail about the specified command\n", command.Usage)

	return nil
}

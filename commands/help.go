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

	commands := make(map[string]CommandHandler)
	for cmdName, handler := range routes.Routes() {
		if cmd, ok := handler.(CommandHandler); ok {
			commands[cmdName] = cmd
		}
	}

	args := request.Args
	if len(args) > 0 {
		cmdName := args[0]
		if cmd, ok := commands[cmdName]; ok {
			cmd.Define(request.FlagSet)
			fmt.Fprintf(writer, "\n%s\n", cmd.GetName())
			fmt.Fprintf(writer, "  Usage: %s\n", cmdName)
			fmt.Fprintf(writer, "  %s\n\n", cmd.GetSummary())
			fmt.Fprintf(writer, "%s\n\n", cmd.GetDescription())

			if usage := request.FlagSet.DefaultUsage(); usage != "" {
				fmt.Fprintln(writer, "\nUsage")
				fmt.Fprintln(writer, usage)
			}

			return nil
		}
	}

	fmt.Fprintf(writer, "\n%s: %s\n", command.Usage, fmt.Sprintf("%s or %s <command-name>", command.Usage, command.Usage))
	fmt.Fprintln(writer, "\nCommands")
	fmt.Fprintln(writer, "------------------")
	for cmdName, cmd := range commands {
		fmt.Fprintf(writer, "%12s:\t%s\n", cmdName, cmd.GetSummary())
	}

	if usage := request.FlagSet.DefaultUsage(); usage != "" {
		fmt.Fprintln(writer, "\nUsage")
		fmt.Fprintln(writer, usage)
	}

	fmt.Fprintf(writer, "\nUse \"%s <command-name>\" for detail about the specified command\n", command.Usage)

	return nil
}

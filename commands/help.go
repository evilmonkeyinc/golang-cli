package commands

import (
	"fmt"
	"sort"

	"github.com/evilmonkeyinc/golang-cli/shell"
)

// HelpCommand allows for basic output of command metadata.
type HelpCommand struct {
	// The command string that should be used to call this function,
	// used in the default function output
	Usage string
}

func (command *HelpCommand) printCommandList(writer shell.ResponseWriter, commands map[string]CommandHandler) {
	if len(commands) > 0 {

		keys := make([]string, 0, len(commands))
		for key := range commands {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		fmt.Fprintln(writer, "\nCommands")
		fmt.Fprintln(writer, "------------------")
		for _, cmdName := range keys {
			cmd := commands[cmdName]
			fmt.Fprintf(writer, "%12s:\t%s\n", cmdName, cmd.GetSummary())
		}
	}
}

func (command *HelpCommand) printCommandHandlerDetails(writer shell.ResponseWriter, request *shell.Request, commandHandler CommandHandler, args []string) error {

	commands := make(map[string]CommandHandler)
	if routes, ok := commandHandler.(shell.Routes); ok {
		for cmdName, handler := range routes.Routes() {
			if cmd, ok := handler.(CommandHandler); ok {
				commands[cmdName] = cmd
			}
		}
	}

	if len(args) > 0 {
		if cmd, ok := commands[args[0]]; ok {
			return command.printCommandHandlerDetails(writer, request, cmd, args[1:])
		}
	}

	commandHandler.Define(request.FlagSet)
	fmt.Fprintf(writer, "\n%s\n", commandHandler.GetName())
	fmt.Fprintf(writer, "  Usage: %s\n", commandHandler.GetUsage())
	fmt.Fprintf(writer, "  %s\n\n", commandHandler.GetSummary())
	fmt.Fprintf(writer, "%s\n\n", commandHandler.GetDescription())

	command.printCommandList(writer, commands)

	if usage := request.FlagSet.DefaultUsage(); usage != "" {
		fmt.Fprintln(writer, "\nUsage")
		fmt.Fprintln(writer, usage)
	}

	return nil
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
		if args[0] == command.Usage && len(args) > 1 {
			// if first arg is "help" then we ignore that,
			// this lets us support any command returning
			// HelpRequested, even the help command
			args = args[1:]
		}
		if cmd, ok := commands[args[0]]; ok {
			return command.printCommandHandlerDetails(writer, request, cmd, args[1:])
		}
	}

	fmt.Fprintf(writer, "\n%s: %s\n", command.Usage, fmt.Sprintf("%s or %s <command-name>", command.Usage, command.Usage))
	command.printCommandList(writer, commands)

	if usage := request.FlagSet.DefaultUsage(); usage != "" {
		fmt.Fprintln(writer, "\nUsage")
		fmt.Fprintln(writer, usage)
	}

	fmt.Fprintf(writer, "\nUse \"%s <command-name>\" for detail about the specified command\n", command.Usage)

	return nil
}

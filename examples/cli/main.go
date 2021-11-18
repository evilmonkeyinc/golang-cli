package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/evilmonkeyinc/golang-cli/pkg/commands"
	"github.com/evilmonkeyinc/golang-cli/pkg/middleware"
	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pingCommand := &commands.Command{
		Name:        "Ping",
		Summary:     "Simple ping pong command",
		Description: "Simple command that will output the word pong",
		Flags: func(fd shell.FlagDefiner) {
			fd.String("suffix", "", "")
		},
		Function: func(rw shell.ResponseWriter, r *shell.Request) error {
			message := "pong"

			if suffix := r.FlagValues().GetString("suffix"); suffix != nil {
				message = fmt.Sprintf("%s%s", message, *suffix)
			}

			if toUpper := r.FlagValues().GetBool("toUpper"); toUpper != nil && *toUpper {
				message = strings.ToUpper(message)
			}

			fmt.Fprintln(rw, message)
			return nil
		},
	}

	newShell := new(shell.Shell)
	newShell.Flags(shell.FlagHandlerFunction(func(fd shell.FlagDefiner) {
		fd.Bool("toUpper", false, "")
	}))
	newShell.Use(middleware.Recoverer())
	newShell.Handle("ping", pingCommand)
	newShell.Route("sub", func(r shell.Router) {
		r.Handle("ping", pingCommand)
	})
	newShell.HandleFunction("secret", func(rw shell.ResponseWriter, r *shell.Request) error {
		panic("this command should not be called.")
	})
	newShell.Handle("help", &commands.HelpCommand{})
	newShell.Execute(ctx)
}

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/evilmonkeyinc/golang-cli/commands"
	"github.com/evilmonkeyinc/golang-cli/flags"
	"github.com/evilmonkeyinc/golang-cli/middleware"
	"github.com/evilmonkeyinc/golang-cli/shell"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pingCommand := &commands.Command{
		Name:        "Ping",
		Summary:     "Simple ping pong command",
		Description: "Simple command that will output the word pong",
		Flags: func(fd flags.FlagDefiner) {
			fd.String("suffix", "", "add a suffix to the response")
		},
		Function: func(rw shell.ResponseWriter, r *shell.Request) error {
			message := "pong"

			if suffix, ok := r.FlagValues().GetString("suffix"); ok {
				message = fmt.Sprintf("%s%s", message, suffix)
			}

			if toUpper, ok := r.FlagValues().GetBool("toUpper"); ok && toUpper {
				message = strings.ToUpper(message)
			}

			fmt.Fprintln(rw, message)
			return nil
		},
	}

	newShell := new(shell.Shell)
	newShell.Flags(flags.FlagHandlerFunction(func(fd flags.FlagDefiner) {
		fd.Bool("toUpper", false, "make the response uppercase")
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

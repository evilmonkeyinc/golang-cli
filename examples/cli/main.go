package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/evilmonkeyinc/golang-cli/commands"
	"github.com/evilmonkeyinc/golang-cli/errors"
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
	newShell.Options(shell.OptionHelpHandler(&commands.HelpCommand{Usage: "help"}))
	newShell.Flags(flags.FlagHandlerFunction(func(fd flags.FlagDefiner) {
		fd.Bool("toUpper", false, "make the response uppercase")
	}))
	newShell.Use(middleware.Recoverer())
	newShell.Handle("ping", pingCommand)
	newShell.Handle("users", commands.NewCommandRouter("Users", "Commands for user management", "A series of commands to aid in user management", "users add|delete|list", func(r shell.Router) {
		r.Handle("list", &commands.Command{
			Name:        "List",
			Summary:     "List users",
			Description: "Will list all valid users",
			Usage:       "list",
			Function: func(rw shell.ResponseWriter, r *shell.Request) error {
				return fmt.Errorf("list function called")
			},
		})
		r.Handle("add", &commands.Command{
			Name:        "Add",
			Summary:     "Add user",
			Description: "Will add a new user",
			Usage:       "add email@example.com",
			Function: func(rw shell.ResponseWriter, r *shell.Request) error {
				return fmt.Errorf("add function called")
			},
		})
		r.Handle("delete", &commands.Command{
			Name:        "Delete",
			Summary:     "Delete user",
			Description: "Will delete an existing user",
			Usage:       "delete email@example.com",
			Function: func(rw shell.ResponseWriter, r *shell.Request) error {
				return fmt.Errorf("delete function called")
			},
		})
	}))

	newShell.HandleFunction("secret", func(rw shell.ResponseWriter, r *shell.Request) error {
		panic("this command should not be called.")
	})
	newShell.HandleFunction("help", func(shell.ResponseWriter, *shell.Request) error {
		return errors.HelpRequested("help command")
	})
	newShell.Execute(ctx)
}

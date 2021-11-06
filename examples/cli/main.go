package main

import (
	"context"
	"fmt"

	"github.com/evilmonkeyinc/golang-cli/pkg/commands"
	"github.com/evilmonkeyinc/golang-cli/pkg/middleware"
	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newShell := new(shell.Shell)
	newShell.Use(middleware.Recoverer())
	newShell.Handle("ping", &commands.Command{
		Name:        "Ping",
		Summary:     "Simple ping pong command",
		Description: "Simple command that will output the word pong",
		Function: func(rw shell.ResponseWriter, r *shell.Request) error {
			fmt.Fprintln(rw, "pong")
			return nil
		},
	})
	newShell.HandleFunction("secret", func(rw shell.ResponseWriter, r *shell.Request) error {
		panic("this command should not be called.")
	})
	newShell.Handle("help", &commands.HelpCommand{})
	newShell.Execute(ctx)
}

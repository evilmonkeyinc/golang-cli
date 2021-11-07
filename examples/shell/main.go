package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/evilmonkeyinc/golang-cli/pkg/commands"
	"github.com/evilmonkeyinc/golang-cli/pkg/middleware"
	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newShell := new(shell.Shell)
	newShell.Options(shell.OptionShellPrompt("example>"))
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
	newShell.HandleFunction("exit", func(rw shell.ResponseWriter, r *shell.Request) error {
		os.Exit(1)
		return nil
	})

	go newShell.Start(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signals:
		cancel()
		<-newShell.Closed()
	case <-newShell.Closed():
		cancel()
	}
}

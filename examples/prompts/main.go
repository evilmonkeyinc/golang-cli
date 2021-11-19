/**
Example of using github.com/manifoldco/promptui for stylized user prompts.

promptui will function in command-line and shell mode without any modification.

shell response writer conforms to io.WriteCloser interface so it can be used for prompt
if the output or error writers have been updated using the shell options if required.
**/
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/evilmonkeyinc/golang-cli/shell"
	"github.com/manifoldco/promptui"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newShell := new(shell.Shell)
	newShell.HandleFunction("login", func(rw shell.ResponseWriter, r *shell.Request) error {

		usernamePrompt := promptui.Prompt{
			Label:   "Username",
			Default: "username",
			Stdout:  rw,
		}

		username, err := usernamePrompt.Run()
		if err != nil {
			return err
		}

		passwordPrompt := promptui.Prompt{
			Label:  "Password",
			Mask:   '*',
			Stdout: rw,
		}

		password, err := passwordPrompt.Run()
		if err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("%s:%s", username, password))

		return nil
	})

	if len(os.Args) > 0 {
		newShell.Execute(ctx)
		return
	}

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

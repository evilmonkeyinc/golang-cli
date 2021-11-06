package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/evilmonkeyinc/golang-cli/pkg/middleware"
	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newShell := new(shell.Shell)
	newShell.Use(middleware.Recoverer())
	newShell.HandleFunction("ack", func(rw shell.ResponseWriter, r *shell.Request) error {
		rw.Write([]byte("ACK ACK ACK\n"))
		return nil
	})
	newShell.HandleFunction("panic", func(shell.ResponseWriter, *shell.Request) error {
		panic(fmt.Errorf("PANIC PANIC"))
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

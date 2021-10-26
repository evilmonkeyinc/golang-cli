package shell

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

type Shell struct {
	router       Router
	reader       io.Reader
	outputWriter io.Writer
	errorWriter  io.Writer
	closed       chan struct{}
}

func (shell *Shell) setup() {
	if shell.router == nil {
		shell.router = newRouter()
	}
	if shell.outputWriter == nil {
		shell.outputWriter = os.Stdout
	}
	if shell.errorWriter == nil {
		shell.errorWriter = os.Stderr
	}
	if shell.reader == nil {
		shell.reader = os.Stdin
	}
	if shell.closed == nil {
		shell.closed = make(chan struct{})
	}
}

func (shell *Shell) execute(ctx context.Context, args []string) error {
	request := newRequest(ctx, args, shell.router)
	writer := newWriter(ctx, shell.outputWriter, shell.errorWriter)
	return shell.router.Execute(writer, request)
}

func (shell *Shell) Options(options ...Option) error {
	for _, option := range options {
		if err := option.Apply(shell); err != nil {
			return err
		}
	}
	return nil
}

func (shell *Shell) Use(middleware ...Middleware) {
	shell.setup()
	shell.router.Use(middleware...)
}

func (shell *Shell) Group(fn func(r Router)) Router {
	shell.setup()
	return shell.router.Group(fn)
}

func (shell *Shell) Route(command string, fn func(r Router)) Router {
	shell.setup()
	return shell.router.Route(command, fn)
}

func (shell *Shell) Handle(command string, handler Handler) {
	shell.setup()
	shell.router.Handle(command, handler)
}

func (shell *Shell) HandleFunction(command string, fn HandlerFunction) {
	shell.setup()
	shell.router.Handle(command, fn)
}

func (shell *Shell) Help(handler HandlerFunction) {
	shell.setup()
	shell.router.Help(handler)
}

func (shell *Shell) NotFound(handler HandlerFunction) {
	shell.setup()
	shell.router.NotFound(handler)
}

func (shell *Shell) Execute(ctx context.Context) error {
	shell.setup()
	return shell.execute(ctx, os.Args[1:])
}

func (shell *Shell) Start(ctx context.Context) error {
	shell.setup()
	reader := bufio.NewReader(shell.reader)

	line := make(chan string)
	for {
		// start a goroutine to get input from the user
		go func(ctx context.Context, input chan<- string) {
			for {
				fmt.Fprintf(shell.outputWriter, "%s ", "shell>") // TODO : shell prompt
				line, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintf(shell.errorWriter, "%v\n", err)
					continue
				}

				input <- line
				return
			}
		}(ctx, line)

		// wait for input or cancel
		select {
		case <-ctx.Done():
			close(shell.closed)
			return nil
		case input := <-line:
			input = strings.TrimSpace(input)
			err := shell.execute(ctx, strings.Split(input, " "))
			if err != nil {
				fmt.Fprintf(shell.errorWriter, "%v\n", err)
				return err
			}
		}
	}
}

func (shell *Shell) Closed() chan struct{} {
	return shell.closed
}

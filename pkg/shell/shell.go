package shell

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// Shell exposes the command-line or interactive shell functionality.
//
// The shell
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
	request := NewRequestWithContext(ctx, []string{}, args, shell.router)
	writer := NewWrapperWriter(ctx, shell.outputWriter, shell.errorWriter)
	return shell.router.Execute(writer, request)
}

// Options will apply the supplied options to the shell.
//
// Options should be called before adding middleware, groups, or handlers.
func (shell *Shell) Options(options ...Option) error {
	for _, option := range options {
		if err := option.Apply(shell); err != nil {
			return err
		}
	}
	return nil
}

// Use appends one or more middleware onto the router stack.
func (shell *Shell) Use(middleware ...Middleware) {
	shell.setup()
	shell.router.Use(middleware...)
}

// Group adds a new inline-router to the router stack.
func (shell *Shell) Group(fn func(r Router)) Router {
	shell.setup()
	return shell.router.Group(fn)
}

// Route adds a new sub-router to the router stack, along the specified command path.
func (shell *Shell) Route(command string, fn func(r Router)) Router {
	shell.setup()
	return shell.router.Route(command, fn)
}

// Handle adds a shell handler to the router stack, along the specified command path.
func (shell *Shell) Handle(command string, handler Handler) {
	shell.setup()
	shell.router.Handle(command, handler)
}

// HandleFunction adds a shell handler function to the router stack, along the specified command path.
func (shell *Shell) HandleFunction(command string, fn HandlerFunction) {
	shell.setup()
	shell.router.HandleFunction(command, fn)
}

// NotFound defines a shell handler that will respond if a command path cannot be evaluated.
func (shell *Shell) NotFound(handler Handler) {
	shell.setup()
	shell.router.NotFound(handler)
}

// Execute is used to execute the shell, using os.Args to evaluate which function to execute.
func (shell *Shell) Execute(ctx context.Context) error {
	shell.setup()
	return shell.execute(ctx, os.Args[1:])
}

// Start is used to begin a new shell session.
//
// The interactive shell will read input and evaluate the commands to execute handler functions.
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
				close(shell.closed)
				return err // TODO : option to exit on error for shell
			}
		}
	}
}

// Closed is used to determine if the shell session is closed.
func (shell *Shell) Closed() chan struct{} {
	return shell.closed
}

[![codecov](https://codecov.io/gh/evilmonkeyinc/golang-cli/branch/main/graph/badge.svg?token=4PU85I7J2R)](https://codecov.io/gh/evilmonkeyinc/golang-cli)
[![coverage](https://github.com/evilmonkeyinc/golang-cli/actions/workflows/coverage.yaml/badge.svg?branch=main)](https://github.com/evilmonkeyinc/golang-cli/actions/workflows/other.yaml)

> This library is on the unstable version v0.X.X, which means there is a chance that any minor update may introduce a breaking change. Where I will endeavor to avoid this, care should be taken updating your dependency on this library until the first stable release v1.0.0 at which point any future breaking changes will result in a new major release.

# Golang-CLI

A simple command-line tool-set for Golang modeled after HTTP routing patterns.

## Install

`go get -u github.com/evilmonkeyinc/golang-cli`

## Usage

### Shell

```golang
import "github.com/evilmonkeyinc/golang-cli/pkg/shell"
...
    newShell := new(shell.Shell)
...
```

### Options

Options adds the ability to customize the shell's properties for your project.

```golang
    newShell.Options(shell.OptionOutputWriter(customOutput))
```

> Options should be set before performing any other actions on the shell

### Middleware

Middleware adds the ability to wrap all shell handler functions in additional logic.

```golang
    newShell.Use(middleware.Recoverer())
```

The most common use-case has been included in the library, panic recovery, which 
makes it possible for the shell to recover from a panic without exiting the interactive-shell

Middleware can also be used to manipulate the shell handler functions ResponseWriter and Request to add or manipulate the existing functionality

### Groups

Groups allow you to define a new inline-router to the shell router stack.

```golang
    newShell.Group(func(r shell.Router) {
		...
	})
```

An inline-routers handlers will be evaluated at the same level as the containing router, but makes it possible to define middleware that will be executed for the inline-router handlers only.

### Routes

Routes allow you do define a sub-router which is executed as a shell handler, but include additional
routing instructions to handle sub-commands for the shell.

```golang
    newShell.Route("users", func(r shell.Router) {
	    r.HandleFunction("add", func(rw shell.ResponseWriter, r *shell.Request) error { return nil })
		r.HandleFunction("list", func(rw shell.ResponseWriter, r *shell.Request) error { return nil })
	})
```

The sub commands would then be executed by including spaces between the commands:

```bash
./yourcli users list
```

Routes will also support specific middleware for these sub-commands in the same way as the inline-routers created by Group.

### Handlers

The Handle and HandleFunction functions add shell handlers to the router stack. 

```golang
    newShell.Handle("help", &commands.HelpCommand{})
	newShell.HandleFunction("ping", func(rw shell.ResponseWriter, r *shell.Request) error {
        fmt.Fprintln(rw, "pong")
		return nil
	})
```

```bash
./yourcli ping
pong
```

## Examples

- [CLI Example](examples/cli/main.go)  
Using this library for a standard command-line interface tool.

- [Shell Example](examples/shell/main.go)  
Using this library for an interactive-shell interface.

## References and Inspirations

The following projects were used as references and inspiration for this project 

- [Valdimir Vivien's GOSH](https://github.com/vladimirvivien/gosh)  
- [Go Chi's Chi Router](https://github.com/go-chi/chi)

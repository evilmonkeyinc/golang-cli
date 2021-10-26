package errors

import (
	"errors"
)

var (
	DuplicateCommandFunction error = errors.New("command already been declared")
	CommandNotFound          error = errors.New("command not found")
	UnableToParseCommand     error = errors.New("unable to parse command")
)

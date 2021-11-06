package errors

import (
	"errors"
	"fmt"
)

var (
	errCommandNotFound  error = errors.New("command not found")
	errDuplicateCommand error = errors.New("command has already been declared")
	errOptionIsSet      error = errors.New("option has already been used or shell has already been initialized")
)

// CommandNotFound returns a command not found error
func CommandNotFound(command string) error {
	return fmt.Errorf("'%s' %w", command, errCommandNotFound)
}

// DuplicateCommand returns a duplicate command error
func DuplicateCommand(command string) error {
	return fmt.Errorf("'%s' %w", command, errDuplicateCommand)
}

// OptionIsSet returns an option is ste error
func OptionIsSet(option string) error {
	return fmt.Errorf("'%s' %w", option, errOptionIsSet)
}

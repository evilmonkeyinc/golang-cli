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

func CommandNotFound(command string) error {
	return fmt.Errorf("'%s' %w", command, errCommandNotFound)
}

func DuplicateCommand(command string) error {
	return fmt.Errorf("'%s' %w", command, errDuplicateCommand)
}

func OptionIsSet(option string) error {
	return fmt.Errorf("'%s' %w", option, errOptionIsSet)
}

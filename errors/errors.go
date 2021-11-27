package errors

import (
	"errors"
	"fmt"
)

var (
	errCommandNotFound    error = errors.New("command not found")
	errDuplicateCommand   error = errors.New("command has already been declared")
	errFlagsetParseFailed error = errors.New("flagset parse failed")
	errHelpRequested      error = errors.New("help requested")
	errOptionIsInvalid    error = errors.New("option paramaters are undefined or invalid")
	errOptionIsSet        error = errors.New("option has already been used or shell has already been initialized")
)

// CommandNotFound returns a command not found error
func CommandNotFound(command string) error {
	return fmt.Errorf("'%s' %w", command, errCommandNotFound)
}

// DuplicateCommand returns a duplicate command error
func DuplicateCommand(command string) error {
	return fmt.Errorf("'%s' %w", command, errDuplicateCommand)
}

// FlagsetParseFailed returns a flagset parse failed error
func FlagsetParseFailed(reason string) error {
	return fmt.Errorf("%w %s", errFlagsetParseFailed, reason)
}

// HelpRequested returns a help requested error
func HelpRequested(reason string) error {
	return fmt.Errorf("%w %s", errHelpRequested, reason)
}

// IsHelpRequested determines if the specified error is a help requested error
func IsHelpRequested(err error) bool {
	return errors.Is(err, errHelpRequested)
}

// OptionIsSet returns an option is set error
func OptionIsSet(option string) error {
	return fmt.Errorf("'%s' %w", option, errOptionIsSet)
}

// OptionIsInvalid returns an option is invalid error
func OptionIsInvalid(option string) error {
	return fmt.Errorf("'%s' %w", option, errOptionIsInvalid)
}

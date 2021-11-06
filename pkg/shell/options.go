package shell

import (
	"io"

	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
)

// The OptionFunction type is an adapter to allow the use of ordinary functions as shell options.
type OptionFunction func(shell *Shell) error

// Apply is used to apply the shell options.
func (option OptionFunction) Apply(shell *Shell) error {
	return option(shell)
}

// The Option interface describes a shell option function
type Option interface {
	// Apply is used to apply the shell options.
	Apply(shell *Shell) error
}

// OptionInput shell option allows the user to set the shell input reader.
func OptionInput(reader io.Reader) Option {
	return &inputOption{
		reader: reader,
	}
}

type inputOption struct {
	reader io.Reader
}

func (option *inputOption) Apply(shell *Shell) error {
	if shell.reader != nil {
		return errors.OptionIsSet("Input")
	}
	shell.reader = option.reader
	return nil
}

// OptionOutputWriter shell option allows the user to set the shell output writer.
func OptionOutputWriter(writer io.Writer) Option {
	return &outputWriterOption{
		writer: writer,
	}
}

type outputWriterOption struct {
	writer io.Writer
}

func (option *outputWriterOption) Apply(shell *Shell) error {
	if shell.outputWriter != nil {
		return errors.OptionIsSet("OutputWriter")
	}
	shell.outputWriter = option.writer
	return nil
}

// OptionErrorWriter shell option allows the user to set the shell error writer.
func OptionErrorWriter(writer io.Writer) Option {
	return &errorWriterOption{
		writer: writer,
	}
}

type errorWriterOption struct {
	writer io.Writer
}

func (option *errorWriterOption) Apply(shell *Shell) error {
	if shell.errorWriter != nil {
		return errors.OptionIsSet("ErrorWriter")
	}
	shell.errorWriter = option.writer
	return nil
}

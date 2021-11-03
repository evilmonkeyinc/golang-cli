package shell

import (
	"io"

	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
)

type OptionFunction func(shell *Shell) error

func (option OptionFunction) Apply(shell *Shell) error {
	return option(shell)
}

type Option interface {
	Apply(shell *Shell) error
}

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

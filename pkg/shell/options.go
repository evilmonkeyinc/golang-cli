package shell

import "io"

type Option interface {
	Apply(shell *Shell) error
}

func Input(reader io.Reader) Option {
	return &inputOption{
		reader: reader,
	}
}

type inputOption struct {
	reader io.Reader
}

func (option *inputOption) Apply(shell *Shell) error {
	shell.reader = option.reader
	return nil
}

func OutputWriter(writer io.Writer) Option {
	return &outputWriterOption{
		writer: writer,
	}
}

type outputWriterOption struct {
	writer io.Writer
}

func (option *outputWriterOption) Apply(shell *Shell) error {
	shell.outputWriter = option.writer
	return nil
}

func ErrorWriter(writer io.Writer) Option {
	return &errorWriterOption{
		writer: writer,
	}
}

type errorWriterOption struct {
	writer io.Writer
}

func (option *errorWriterOption) Apply(shell *Shell) error {
	shell.errorWriter = option.writer
	return nil
}

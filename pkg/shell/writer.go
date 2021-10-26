package shell

import (
	"context"
	"io"
)

type ResponseWriter interface {
	Write([]byte) (int, error)
	WriteError([]byte) (int, error)
}

func newWriter(ctx context.Context, outputWritter io.Writer, errorWriter io.Writer) *writer {
	return &writer{
		outputWritter: outputWritter,
		errorWritter:  errorWriter,
	}
}

type writer struct {
	outputWritter io.Writer
	errorWritter  io.Writer
}

func (writer *writer) Write(bytes []byte) (int, error) {
	return writer.outputWritter.Write(bytes)
}

func (writer *writer) WriteError(bytes []byte) (int, error) {
	return writer.errorWritter.Write(bytes)
}

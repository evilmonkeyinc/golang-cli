package shell

import (
	"context"
	"io"
	"os"
)

// ResponseWriter is an interface is used by an shell handler to construct a response.
type ResponseWriter interface {
	// Close will close the output and error writers if the conform to io.WriteCloser.
	Close() error
	// ErrorWriter returns an io.Writer used for the error output
	ErrorWriter() io.Writer
	// Write byte array to the output writer.
	Write([]byte) (int, error)
	// WriteError will write byte array to the error writer.
	WriteError([]byte) (int, error)
}

// NewWrapperWriter returns a new WrapperWriter using the specified output and error writers.
func NewWrapperWriter(ctx context.Context, outputWritter io.Writer, errorWriter io.Writer) *WrapperWriter {
	return &WrapperWriter{
		outputWritter: outputWritter,
		errorWritter:  errorWriter,
	}
}

// The WrapperWriter is used as a simple ResponseWritter
type WrapperWriter struct {
	outputWritter io.Writer
	errorWritter  io.Writer
}

// Close will close the output and error writers if the conform to io.WriteCloser.
func (writer *WrapperWriter) Close() error {
	if closer, ok := writer.errorWritter.(io.WriteCloser); ok {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	if closer, ok := writer.outputWritter.(io.WriteCloser); ok {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// ErrorWriter returns an io.Writer used for the error output
func (writer *WrapperWriter) ErrorWriter() io.Writer {
	if writer.errorWritter == nil {
		writer.errorWritter = os.Stderr
	}
	return writer.errorWritter
}

// Write byte array to the output writer.
func (writer *WrapperWriter) Write(bytes []byte) (int, error) {
	if writer.outputWritter == nil {
		writer.outputWritter = os.Stdout
	}
	return writer.outputWritter.Write(bytes)
}

// WriteError will write byte array to the error writer.
func (writer *WrapperWriter) WriteError(bytes []byte) (int, error) {
	if writer.errorWritter == nil {
		writer.errorWritter = os.Stderr
	}
	return writer.errorWritter.Write(bytes)
}

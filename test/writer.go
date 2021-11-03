package test

import "fmt"

type TestResponseWriter struct {
	WriteData string
	ErrorData string
}

func (writer *TestResponseWriter) Write(p []byte) (n int, err error) {
	writer.WriteData = fmt.Sprintf("%s%s", writer.WriteData, string(p))
	return len(writer.WriteData), nil
}

func (writer *TestResponseWriter) WriteError(p []byte) (n int, err error) {
	writer.ErrorData = fmt.Sprintf("%s%s", writer.ErrorData, string(p))
	return len(writer.ErrorData), nil
}

package shell

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testWriter implements io.Writer so it can be used to test ResponseWriter
type testWriter struct {
	Prefix string
	Data   string
}

func (writer *testWriter) Write(p []byte) (n int, err error) {
	writer.Data = fmt.Sprintf("%s%s%s", writer.Data, writer.Prefix, string(p))
	return len(writer.Data), nil
}

// Validate the writer struct matches the ResponseWriter interface
var _ ResponseWriter = &writer{}

func Test_ResponseWriter(t *testing.T) {

	type input struct {
		writeInput string
		errorInput string
		prefix     string
	}

	type expected struct {
		writeOutput string
		errorOutput string
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "test 1",
			input: input{
				prefix:     "",
				writeInput: "input",
				errorInput: "error",
			},
			expected: expected{
				writeOutput: "input",
				errorOutput: "error",
			},
		},
		{
			name: "test 2",
			input: input{
				prefix:     "prefix",
				writeInput: "",
				errorInput: "error",
			},
			expected: expected{
				writeOutput: "prefix",
				errorOutput: "prefixerror",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			outputWriter := &testWriter{
				Prefix: test.input.prefix,
			}
			errorWriter := &testWriter{
				Prefix: test.input.prefix,
			}

			writer := newWriter(context.Background(), outputWriter, errorWriter)
			writer.Write([]byte(test.input.writeInput))
			writer.WriteError([]byte(test.input.errorInput))

			assert.EqualValues(t, test.expected.writeOutput, outputWriter.Data)
			assert.EqualValues(t, test.expected.errorOutput, errorWriter.Data)
		})
	}

}

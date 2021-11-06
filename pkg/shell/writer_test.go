package shell

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Validate the writer struct matches the ResponseWriter interface
var _ ResponseWriter = &WrapperWriter{}

func Test_WrapperWriter_defaults(t *testing.T) {

	writer := &WrapperWriter{}

	writer.Write([]byte(""))
	writer.WriteError([]byte(""))

	assert.Equal(t, os.Stdout, writer.outputWritter)
	assert.Equal(t, os.Stderr, writer.errorWritter)

	writer.errorWritter = nil
	writer.ErrorWriter()
	assert.Equal(t, os.Stderr, writer.errorWritter)
}

func Test_WrapperWriter(t *testing.T) {

	type input struct {
		writeInput string
		errorInput string
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
				writeInput: "",
				errorInput: "error",
			},
			expected: expected{
				writeOutput: "",
				errorOutput: "error",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			outputWriter := &bytes.Buffer{}
			errorWriter := &bytes.Buffer{}

			writer := NewWrapperWriter(context.Background(), outputWriter, errorWriter)
			writer.Write([]byte(test.input.writeInput))
			writer.WriteError([]byte(test.input.errorInput))

			assert.EqualValues(t, test.expected.writeOutput, outputWriter.String())
			assert.EqualValues(t, test.expected.errorOutput, errorWriter.String())

			assert.Equal(t, errorWriter, writer.ErrorWriter())
		})
	}
}

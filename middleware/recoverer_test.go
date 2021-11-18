package middleware

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/shell"
	"github.com/stretchr/testify/assert"
)

func Test_Recoverer(t *testing.T) {

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string",
			input:    "this panic should be logged",
			expected: "this panic should be logged",
		},
		{
			name:     "error",
			input:    fmt.Errorf("this panic should be logged"),
			expected: "this panic should be logged",
		},
		{
			name:     "map",
			input:    map[string]string{"error": "object logged"},
			expected: "map[error:object logged]",
		},
	}

	testRequest := shell.NewRequest([]string{}, []string{}, &shell.DefaultFlagSet{}, nil)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errorWriter := &bytes.Buffer{}
			testWritter := shell.NewWrapperWriter(context.Background(), &bytes.Buffer{}, errorWriter)

			next := shell.HandlerFunction(func(rw shell.ResponseWriter, r *shell.Request) error {
				panic(test.input)
			})

			middleware := Recoverer()
			middleware.Handle(next).Execute(testWritter, testRequest.WithContext(context.Background()))

			actual := strings.Split(errorWriter.String(), "\n")
			firstLine := actual[0]
			assert.Equal(t, test.expected, firstLine)
		})
	}
}

package middleware

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
	"github.com/evilmonkeyinc/golang-cli/test"
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

	testWritter := &test.TestResponseWriter{}
	testRequest := &test.TestRequest{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testWritter.ErrorData = ""

			next := shell.HandlerFunction(func(rw shell.ResponseWriter, r shell.Request) error {
				panic(test.input)
			})

			middleware := Recoverer()
			middleware.Handle(next).Execute(testWritter, testRequest.WithContext(context.Background()))
			firstLine := strings.Split(testWritter.ErrorData, "\n")[0]
			assert.Equal(t, test.expected, firstLine)
		})
	}
}

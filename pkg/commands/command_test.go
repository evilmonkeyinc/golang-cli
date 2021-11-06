package commands

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
	"github.com/stretchr/testify/assert"
)

func Test_Command(t *testing.T) {

	tests := []struct {
		name     string
		input    *Command
		expected error
	}{
		{
			name: "no function",
			input: &Command{
				Name: "missing",
			},
			expected: errors.CommandNotFound("missing"),
		},
		{
			name: "found",
			input: &Command{
				Name: "found",
				Function: func(shell.ResponseWriter, *shell.Request) error {
					return fmt.Errorf("found")
				},
			},
			expected: fmt.Errorf("found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.input.Execute(nil, nil)
			assert.Equal(t, test.expected, actual)
		})
	}

}

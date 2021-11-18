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
		{
			name: "found with flags",
			input: &Command{
				Name: "found",
				Flags: func(fd shell.FlagDefiner) {
					fd.Bool("suffix", false, "")
				},
				Function: func(rw shell.ResponseWriter, r *shell.Request) error {
					includeSuffix := r.FlagValues().GetBool("suffix")
					if includeSuffix != nil && *includeSuffix {
						return fmt.Errorf("found with suffix")
					}
					return fmt.Errorf("found without suffix")
				},
			},
			expected: fmt.Errorf("found with suffix"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flagSet := &shell.DefaultFlagSet{}
			test.input.Define(flagSet)
			flagSet.Parse([]string{"-suffix"})
			actual := test.input.Execute(nil, shell.NewRequest(nil, nil, flagSet, nil))
			assert.Equal(t, test.expected, actual)
		})
	}

}

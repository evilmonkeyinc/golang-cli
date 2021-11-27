package commands

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/errors"
	"github.com/evilmonkeyinc/golang-cli/flags"
	"github.com/evilmonkeyinc/golang-cli/shell"
	"github.com/stretchr/testify/assert"
)

// Validate the Command struct matches the CommandHandler interface
var _ CommandHandler = &Command{}

func Test_Command(t *testing.T) {

	t.Run("Getters", func(t *testing.T) {
		command := &Command{
			Name:        "Name",
			Summary:     "The command summary",
			Description: "The command description",
			Usage:       "name <arg1>",
		}
		assert.Equal(t, "Name", command.GetName())
		assert.Equal(t, "The command summary", command.GetSummary())
		assert.Equal(t, "The command description", command.GetDescription())
		assert.Equal(t, "name <arg1>", command.GetUsage())
	})

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
				Flags: func(fd flags.FlagDefiner) {
					fd.Bool("suffix", false, "")
				},
				Function: func(rw shell.ResponseWriter, r *shell.Request) error {
					includeSuffix, ok := r.FlagValues().GetBool("suffix")
					if ok && includeSuffix {
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
			flagSet := &flags.DefaultFlagSet{}
			test.input.Define(flagSet)
			flagSet.Parse([]string{"-suffix"})
			actual := test.input.Execute(nil, shell.NewRequest(nil, nil, flagSet, nil))
			assert.Equal(t, test.expected, actual)
		})
	}

}

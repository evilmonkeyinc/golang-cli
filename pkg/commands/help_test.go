package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/pkg/middleware"
	"github.com/evilmonkeyinc/golang-cli/pkg/shell"
	"github.com/stretchr/testify/assert"
)

func Test_HelpCommand(t *testing.T) {

	tests := []struct {
		name     string
		input    []string
		usage    string
		expected []string
	}{
		{
			name:     "ping",
			input:    []string{"ping"},
			usage:    "help",
			expected: []string{"pong"},
		},
		{
			name:  "help",
			input: []string{"help"},
			usage: "help",
			expected: []string{
				"",
				"help: help or help <command-name>",
				"",
				"Available commands",
				"------------------",
				"        ping:	Simple ping pong command",
				"",
				"Use \"help <command-name>\" for detail about the specified command",
			},
		},
		{
			name:  "help no usage",
			input: []string{"help"},
			usage: "",
			expected: []string{
				"",
				"help: help or help <command-name>",
				"",
				"Available commands",
				"------------------",
				"        ping:	Simple ping pong command",
				"",
				"Use \"help <command-name>\" for detail about the specified command",
			},
		},
		{
			name:  "help ping",
			input: []string{"help", "ping"},
			usage: "help",
			expected: []string{
				"",
				"ping",
				"  Usage: Ping",
				"  Simple ping pong command",
				"",
				"Simple command that will output the word pong",
				"",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testWriter := &bytes.Buffer{}

			newShell := new(shell.Shell)
			newShell.Options(shell.OptionOutputWriter(testWriter))
			newShell.Use(middleware.Recoverer())
			newShell.Handle("ping", &Command{
				Name:        "Ping",
				Summary:     "Simple ping pong command",
				Description: "Simple command that will output the word pong",
				Function: func(rw shell.ResponseWriter, r *shell.Request) error {
					fmt.Fprintln(rw, "pong")
					return nil
				},
			})
			newShell.HandleFunction("secret", func(rw shell.ResponseWriter, r *shell.Request) error {
				panic("this command should not be called.")
			})
			newShell.Handle("help", &HelpCommand{Usage: test.usage})

			os.Args = append([]string{"cmd"}, test.input...)
			err := newShell.Execute(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, strings.Join(test.expected, "\n")+"\n", testWriter.String())
		})
	}

}

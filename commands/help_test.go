package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/flags"
	"github.com/evilmonkeyinc/golang-cli/middleware"
	"github.com/evilmonkeyinc/golang-cli/shell"
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
				"Commands",
				"------------------",
				"        ping:\tSimple ping pong command",
				"",
				"Usage",
				"  -toUpper",
				"    \tstate if the response should be uppercase",
				"",
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
				"Commands",
				"------------------",
				"        ping:	Simple ping pong command",
				"",
				"Usage",
				"  -toUpper",
				"    \tstate if the response should be uppercase",
				"",
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
				"Ping",
				"  Usage: ping",
				"  Simple ping pong command",
				"",
				"Simple command that will output the word pong",
				"",
				"",
				"Usage",
				"  -suffix string",
				"    \ta suffix for the function response",
				"  -toUpper",
				"    \tstate if the response should be uppercase",
				"",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testWriter := &bytes.Buffer{}

			newShell := new(shell.Shell)
			newShell.Options(shell.OptionOutputWriter(testWriter))
			newShell.Flags(flags.FlagHandlerFunction(func(fd flags.FlagDefiner) {
				fd.Bool("toUpper", false, "state if the response should be uppercase")
			}))
			newShell.Use(middleware.Recoverer())
			newShell.Handle("ping", &Command{
				Name:        "Ping",
				Summary:     "Simple ping pong command",
				Description: "Simple command that will output the word pong",
				Flags: func(fd flags.FlagDefiner) {
					fd.String("suffix", "", "a suffix for the function response")
				},
				Function: func(rw shell.ResponseWriter, r *shell.Request) error {
					suffix, _ := r.FlagSet.GetString("suffix")

					response := fmt.Sprintf("pong%s", suffix)
					if toUpper, _ := r.FlagSet.GetBool("toUpper"); toUpper {
						response = strings.ToUpper(response)
					}

					fmt.Fprintln(rw, response)
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

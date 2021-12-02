package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/errors"
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
				"       users:\tCommands for user management",
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
				"Commands",
				"------------------",
				"        ping:\tSimple ping pong command",
				"       users:\tCommands for user management",
				"",
				"Usage",
				"  -toUpper",
				"    \tstate if the response should be uppercase",
				"",
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
		{
			name:  "help users",
			input: []string{"help", "users"},
			usage: "",
			expected: []string{
				"",
				"Users",
				"  Usage: users list|add|delete",
				"  Commands for user management",
				"",
				"A series of commands to aid in user management",
				"",
				"",
				"Commands",
				"------------------",
				"         add:	Add user",
				"      delete:	Delete user",
				"        list:	List users",
				"",
				"Usage",
				"  -toUpper",
				"    	state if the response should be uppercase",
				"",
			},
		},
		{
			name:  "help users add",
			input: []string{"help", "users", "add"},
			usage: "",
			expected: []string{
				"",
				"Add",
				"  Usage: add email@example.com",
				"  Add user",
				"",
				"Will add a new user", "", "",
				"Usage",
				"  -toUpper",
				"    	state if the response should be uppercase",
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
				Usage:       "ping",
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
			newShell.Handle("users", NewCommandRouter("Users", "Commands for user management", "A series of commands to aid in user management", "users list|add|delete", func(r shell.Router) {
				r.Handle("list", &Command{
					Name:        "List",
					Summary:     "List users",
					Description: "Will list all valid users",
					Usage:       "list",
					Function: func(rw shell.ResponseWriter, r *shell.Request) error {
						return fmt.Errorf("list function called")
					},
				})
				r.Handle("add", &Command{
					Name:        "Add",
					Summary:     "Add user",
					Description: "Will add a new user",
					Usage:       "add email@example.com",
					Function: func(rw shell.ResponseWriter, r *shell.Request) error {
						return fmt.Errorf("add function called")
					},
				})
				r.Handle("delete", &Command{
					Name:        "Delete",
					Summary:     "Delete user",
					Description: "Will delete an existing user",
					Usage:       "delete email@example.com",
					Function: func(rw shell.ResponseWriter, r *shell.Request) error {
						return fmt.Errorf("delete function called")
					},
				})
			}))
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

func Test_HelpCommandOption(t *testing.T) {

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "base",
			input:    []string{"ping"},
			expected: []string{"pong"},
		},
		{
			name:  "command",
			input: []string{"help"},
			expected: []string{
				"",
				"help: help or help <command-name>",
				"",
				"Commands",
				"------------------",
				"        ping:\tPing/Pong",
				"",
				"Use \"help <command-name>\" for detail about the specified command",
			},
		},
		{
			name:  "flag",
			input: []string{"-help"},
			expected: []string{
				"",
				"help: help or help <command-name>",
				"",
				"Commands",
				"------------------",
				"        ping:\tPing/Pong",
				"",
				"Use \"help <command-name>\" for detail about the specified command",
			},
		},
		{
			name:  "help ping",
			input: []string{"help", "ping"},
			expected: []string{
				"",
				"Ping",
				"  Usage: ping",
				"  Ping/Pong",
				"",
				"This is the description string",
				"",
				"",
				"Usage",
				"  -deprecated",
				"    	do not use this flag",
				"",
			},
		},
		{
			name:  "flag after",
			input: []string{"ping", "-h"},
			expected: []string{
				"",
				"Ping",
				"  Usage: ping",
				"  Ping/Pong",
				"",
				"This is the description string",
				"",
				"",
				"Usage",
				"  -deprecated",
				"    	do not use this flag",
				"",
			},
		},
		{
			name:  "flag before",
			input: []string{"-h", "ping"},
			expected: []string{
				"",
				"Ping",
				"  Usage: ping",
				"  Ping/Pong",
				"",
				"This is the description string",
				"",
				"",
				"Usage",
				"  -deprecated",
				"    	do not use this flag",
				"",
			},
		},
		{
			name:  "ping requests help",
			input: []string{"ping", "-deprecated"},
			expected: []string{
				"",
				"Ping",
				"  Usage: ping",
				"  Ping/Pong",
				"",
				"This is the description string",
				"",
				"",
				"Usage",
				"  -deprecated",
				"    	do not use this flag",
				"",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			testWriter := &bytes.Buffer{}
			newShell := &shell.Shell{}
			newShell.Options(
				shell.OptionOutputWriter(testWriter),
				shell.OptionHelpHandler(&HelpCommand{Usage: "help"}),
			)
			newShell.HandleFunction("help", func(shell.ResponseWriter, *shell.Request) error {
				return errors.HelpRequested("help function")
			})
			newShell.Handle("ping", &Command{
				Name:        "Ping",
				Summary:     "Ping/Pong",
				Description: "This is the description string",
				Usage:       "ping",
				Flags: func(fd flags.FlagDefiner) {
					fd.Bool("deprecated", false, "do not use this flag")
				},
				Function: func(rw shell.ResponseWriter, r *shell.Request) error {
					if deprecated, _ := r.FlagSet.GetBool("deprecated"); deprecated {
						return errors.HelpRequested("bad flag")
					}
					_, err := fmt.Fprintln(rw, "pong")
					return err
				},
			})

			os.Args = append([]string{"cmd"}, test.input...)
			err := newShell.Execute(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, strings.Join(test.expected, "\n")+"\n", testWriter.String())
		})
	}
}

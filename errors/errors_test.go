package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CommandNotFound(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard",
			input:    "PING",
			expected: "'PING' command not found",
		},
		{
			name:     "empty",
			input:    "",
			expected: "'' command not found",
		},
		{
			name:     "lowercase",
			input:    "ping",
			expected: "'ping' command not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := CommandNotFound(test.input)
			assert.Equal(t, test.expected, actual.Error())
			assert.True(t, errors.Is(actual, errCommandNotFound))
		})
	}
}

func Test_DuplicateCommand(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard",
			input:    "PING",
			expected: "'PING' command has already been declared",
		},
		{
			name:     "empty",
			input:    "",
			expected: "'' command has already been declared",
		},
		{
			name:     "lowercase",
			input:    "ping",
			expected: "'ping' command has already been declared",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := DuplicateCommand(test.input)
			assert.Equal(t, test.expected, actual.Error())
			assert.True(t, errors.Is(actual, errDuplicateCommand))
		})
	}
}

func Test_FlagsetParseFailed(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard",
			input:    "PING",
			expected: "flagset parse failed PING",
		},
		{
			name:     "empty",
			input:    "",
			expected: "flagset parse failed ",
		},
		{
			name:     "lowercase",
			input:    "ping",
			expected: "flagset parse failed ping",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := FlagsetParseFailed(test.input)
			assert.Equal(t, test.expected, actual.Error())
			assert.True(t, errors.Is(actual, errFlagsetParseFailed))
		})
	}
}

func Test_HelpRequested(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard",
			input:    "PING",
			expected: "help requested PING",
		},
		{
			name:     "empty",
			input:    "",
			expected: "help requested ",
		},
		{
			name:     "lowercase",
			input:    "ping",
			expected: "help requested ping",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := HelpRequested(test.input)
			assert.Equal(t, test.expected, actual.Error())
			assert.True(t, errors.Is(actual, errHelpRequested))
		})
	}
}

func Test_IsHelpRequested(t *testing.T) {

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "error",
			err:      fmt.Errorf("help"),
			expected: false,
		},
		{
			name:     "help",
			err:      HelpRequested("me"),
			expected: true,
		},
		{
			name:     "wrongly wrapped",
			err:      fmt.Errorf("%v", HelpRequested("me")),
			expected: false,
		},
		{
			name:     "wrapped",
			err:      fmt.Errorf("%w", HelpRequested("me")),
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := IsHelpRequested(test.err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_OptionIsSet(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard",
			input:    "Type",
			expected: "'Type' option has already been used or shell has already been initialized",
		},
		{
			name:     "empty",
			input:    "",
			expected: "'' option has already been used or shell has already been initialized",
		},
		{
			name:     "lowercase",
			input:    "type",
			expected: "'type' option has already been used or shell has already been initialized",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := OptionIsSet(test.input)
			assert.Equal(t, test.expected, actual.Error())
			assert.True(t, errors.Is(actual, errOptionIsSet))
		})
	}
}
func Test_OptionIsInvalid(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard",
			input:    "Type",
			expected: "'Type' option paramaters are undefined or invalid",
		},
		{
			name:     "empty",
			input:    "",
			expected: "'' option paramaters are undefined or invalid",
		},
		{
			name:     "lowercase",
			input:    "type",
			expected: "'type' option paramaters are undefined or invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := OptionIsInvalid(test.input)
			assert.Equal(t, test.expected, actual.Error())
			assert.True(t, errors.Is(actual, errOptionIsInvalid))
		})
	}
}

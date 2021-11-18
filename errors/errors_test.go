package errors

import (
	"errors"
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

package shell

import (
	"context"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/flags"
	"github.com/stretchr/testify/assert"
)

func Test_NewRequest(t *testing.T) {
	path := []string{}
	args := []string{"args"}
	flagSet := &flags.DefaultFlagSet{}

	actual := NewRequest(path, args, flagSet, nil)
	assert.Equal(t, args, actual.Args)
	assert.Equal(t, path, actual.Path)
	assert.Equal(t, nil, actual.Routes)
	assert.Equal(t, flagSet, actual.FlagSet)
}

func Test_NewRequestWithContext(t *testing.T) {

	ctx := context.Background()
	path := []string{}
	args := []string{"args"}
	flagSet := &flags.DefaultFlagSet{}

	actual := NewRequestWithContext(ctx, path, args, flagSet, nil)
	assert.Equal(t, actual.Args, args)
	assert.Equal(t, actual.Path, path)
	assert.Equal(t, actual.ctx, ctx)
	assert.Equal(t, actual.FlagSet, flagSet)
	assert.Equal(t, actual.Routes, nil)
}

func Test_Request_Context(t *testing.T) {
	ctx := context.Background()
	args := []string{"args"}
	path := []string{}

	actual := NewRequestWithContext(ctx, path, args, &flags.DefaultFlagSet{}, nil)
	assert.Equal(t, actual.Context(), ctx)
}

func Test_Request_WithContext(t *testing.T) {
	ctx := context.Background()
	args := []string{"args"}
	path := []string{}

	actual := NewRequestWithContext(ctx, path, args, &flags.DefaultFlagSet{}, nil)
	assert.Equal(t, actual.Context(), ctx)

	type key string
	var ctxKey key = "key"
	nextCtx := context.WithValue(ctx, ctxKey, "value")
	actual = actual.WithContext(nextCtx)
	assert.Equal(t, actual.Context(), nextCtx)
}

func Test_Request_FlagValues(t *testing.T) {

	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]interface{}
	}{
		{
			name:  "none",
			input: map[string]string{},
			expected: map[string]interface{}{
				"bool":   false,
				"string": "",
			},
		},
		{
			name: "all",
			input: map[string]string{
				"string": "a value",
				"bool":   "true",
			},
			expected: map[string]interface{}{
				"bool":   true,
				"string": "a value",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			flagSet := flags.NewDefaultFlagSet()
			flagSet.String("string", "", "")
			flagSet.Bool("bool", false, "")

			for key, value := range test.input {
				flagSet.Set(key, value)
			}

			request := NewRequest(nil, nil, flagSet, nil)

			flagValues := request.FlagValues()
			for key, expected := range test.expected {
				actual := flagValues.Get(key)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func Test_Request_UpdateRequest(t *testing.T) {
	type input struct {
		selectedRoute string
	}

	type expected struct {
		path []string
		args []string
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "expected",
			input: input{
				selectedRoute: "one",
			},
			expected: expected{
				path: []string{"one"},
				args: []string{"two", "three", "four"},
			},
		},
		{
			name: "empty",
			input: input{
				selectedRoute: "",
			},
			expected: expected{
				path: []string{},
				args: []string{"one", "two", "three", "four"},
			},
		},
		{
			name: "unknown",
			input: input{
				selectedRoute: "unknown",
			},
			expected: expected{
				path: []string{"unknown"},
				args: []string{"one", "two", "three", "four"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			args := []string{"one", "two", "three", "four"}
			path := []string{}

			original := NewRequestWithContext(ctx, path, args, &flags.DefaultFlagSet{}, nil)
			updated := original.UpdateRequest(test.input.selectedRoute, nil, nil, nil)

			assert.Equal(t, ctx, updated.ctx)
			assert.Equal(t, test.expected.args, updated.Args)
			assert.Equal(t, test.expected.path, updated.Path)
		})
	}
}

package shell

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRequest(t *testing.T) {
	path := []string{}
	args := []string{"args"}

	actual := NewRequest(path, args, nil)
	assert.Equal(t, actual.Args, args)
	assert.Equal(t, actual.Path, path)
	assert.Equal(t, actual.Routes, nil)
}

func Test_NewRequestWithContext(t *testing.T) {

	ctx := context.Background()
	path := []string{}
	args := []string{"args"}

	actual := NewRequestWithContext(ctx, path, args, nil)
	assert.Equal(t, actual.Args, args)
	assert.Equal(t, actual.Path, path)
	assert.Equal(t, actual.ctx, ctx)
	assert.Equal(t, actual.Routes, nil)
}

func Test_Request_Context(t *testing.T) {
	ctx := context.Background()
	args := []string{"args"}
	path := []string{}

	actual := NewRequestWithContext(ctx, path, args, nil)
	assert.Equal(t, actual.Context(), ctx)
}

func Test_Request_WithContext(t *testing.T) {
	ctx := context.Background()
	args := []string{"args"}
	path := []string{}

	actual := NewRequestWithContext(ctx, path, args, nil)
	assert.Equal(t, actual.Context(), ctx)

	type key string
	var ctxKey key = "key"
	nextCtx := context.WithValue(ctx, ctxKey, "value")
	actual = actual.WithContext(nextCtx)
	assert.Equal(t, actual.Context(), nextCtx)
}

func Test_Request_WithRoutes(t *testing.T) {
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

			original := NewRequestWithContext(ctx, path, args, nil)
			updated := original.WithRoutes(test.input.selectedRoute, nil)

			assert.Equal(t, ctx, updated.ctx)
			assert.Equal(t, test.expected.args, updated.Args)
			assert.Equal(t, test.expected.path, updated.Path)
		})
	}
}

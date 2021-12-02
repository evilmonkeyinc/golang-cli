package middleware

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/shell"
	"github.com/stretchr/testify/assert"
)

func Test_ErrorHandler(t *testing.T) {

	target := fmt.Errorf("target")
	caught := fmt.Errorf("caught")

	tests := []struct {
		name     string
		input    error
		expected error
	}{
		{
			name:     "no error",
			input:    nil,
			expected: nil,
		},
		{
			name:     "unexpected error",
			input:    fmt.Errorf("fallthrough"),
			expected: fmt.Errorf("fallthrough"),
		},
		{
			name:     "expected error",
			input:    target,
			expected: caught,
		},
		{
			name:     "expected error wrapped",
			input:    fmt.Errorf("%w wrapped", target),
			expected: caught,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			middleware := ErrorHandler(func(e error) bool {
				return errors.Is(e, target)
			}, shell.HandlerFunction(func(rw shell.ResponseWriter, r *shell.Request) error {
				return caught
			}))

			request := shell.NewRequest(nil, nil, nil, nil)

			actual := middleware.Handle(shell.HandlerFunction(func(rw shell.ResponseWriter, r *shell.Request) error {
				return test.input
			})).Execute(nil, request)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_ErrorHandlerCaughtError(t *testing.T) {

	t.Run("not set", func(t *testing.T) {
		ctx := context.Background()
		err := ErrorHandlerCaughtError(ctx)
		assert.Nil(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, ctxValueError, nil)
		err := ErrorHandlerCaughtError(ctx)
		assert.Nil(t, err)
	})

	t.Run("string", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, ctxValueError, "not an error")
		err := ErrorHandlerCaughtError(ctx)
		assert.Nil(t, err)
	})
	t.Run("error", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, ctxValueError, fmt.Errorf("valid error"))
		err := ErrorHandlerCaughtError(ctx)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "valid error")
	})
}

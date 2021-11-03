package shell

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Validate the request struct matches the Request interface
var _ Request = &request{}

func Test_Request(t *testing.T) {

	t.Run("newRequest", func(t *testing.T) {
		ctx := context.Background()
		args := []string{"args"}

		actual := newRequest(ctx, args, nil)
		assert.Equal(t, actual.args, args)
		assert.Equal(t, actual.ctx, ctx)
		assert.Equal(t, actual.routes, nil)
	})

	t.Run("Args", func(t *testing.T) {
		ctx := context.Background()
		args := []string{"args"}

		actual := newRequest(ctx, args, nil)
		assert.Equal(t, actual.Args(), args)
	})

	t.Run("Context", func(t *testing.T) {
		ctx := context.Background()
		args := []string{"args"}

		actual := newRequest(ctx, args, nil)
		assert.Equal(t, actual.Context(), ctx)
	})

	t.Run("Routes", func(t *testing.T) {
		ctx := context.Background()
		args := []string{"args"}

		actual := newRequest(ctx, args, nil)
		assert.Equal(t, actual.Routes(), nil)
	})

	t.Run("WithContext", func(t *testing.T) {
		ctx := context.Background()
		args := []string{"args"}

		var actual Request = newRequest(ctx, args, nil)
		assert.Equal(t, actual.Context(), ctx)

		type key string
		var ctxKey key = "key"
		nextCtx := context.WithValue(ctx, ctxKey, "value")
		actual = actual.WithContext(nextCtx)
		assert.Equal(t, actual.Context(), nextCtx)
	})

	t.Run("WithArgs", func(t *testing.T) {
		ctx := context.Background()
		args := []string{"args"}

		var actual Request = newRequest(ctx, args, nil)
		assert.Equal(t, actual.Args(), args)

		updatedArgs := []string{"new"}
		actual = actual.WithArgs(updatedArgs)
		assert.Equal(t, actual.Args(), updatedArgs)
	})

	t.Run("WithRoutes", func(t *testing.T) {
		ctx := context.Background()
		args := []string{"args"}

		var actual Request = newRequest(ctx, args, nil)
		assert.Equal(t, actual.Routes(), nil)

		router := newRouter()
		actual = actual.WithRoutes(router)
		assert.Equal(t, actual.Routes(), router)
	})

}

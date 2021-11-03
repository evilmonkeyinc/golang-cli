package shell

import (
	"context"
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// Validate the router struct matches the Router interface
var _ Router = &router{}

func Test_Router(t *testing.T) {

	t.Run("newRouter", func(t *testing.T) {
		actual := newRouter()
		assert.NotNil(t, actual.children)
		assert.NotNil(t, actual.handlers)
		assert.NotNil(t, actual.middleware)
		assert.Nil(t, actual.parent)
		assert.Nil(t, actual.helpHander)
		assert.Nil(t, actual.notFoundHandler)
	})

	t.Run("childRouter", func(t *testing.T) {
		input := newRouter()
		input.helpHander = HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("help")
		})
		input.notFoundHandler = HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("not found")
		})
		input.handlers["test"] = HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("test")
		})
		input.middleware = append(input.middleware, MiddlewareFunction(func(next Handler) Handler { return next }))
		input.children = append(input.children, newRouter())

		actual := childRouter(input)
		assert.Equal(t, input, actual.parent)
		assert.NotNil(t, actual.helpHander)
		assert.NotNil(t, actual.notFoundHandler)

		assert.NotEqual(t, input.handlers, actual.handlers)
		assert.NotEqual(t, input.middleware, actual.middleware)
		assert.NotEqual(t, input.children, actual.children)
	})

	t.Run("subRouter", func(t *testing.T) {
		input := newRouter()
		input.helpHander = HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("help")
		})
		input.notFoundHandler = HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("not found")
		})
		input.handlers["test"] = HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("test")
		})
		input.middleware = append(input.middleware, MiddlewareFunction(func(next Handler) Handler { return next }))
		input.children = append(input.children, newRouter())

		actual := subRouter(input)
		assert.Nil(t, actual.parent)
		assert.NotNil(t, actual.helpHander)
		assert.NotNil(t, actual.notFoundHandler)

		assert.NotEqual(t, input.handlers, actual.handlers)
		assert.NotEqual(t, input.middleware, actual.middleware)
		assert.NotEqual(t, input.children, actual.children)
	})

	t.Run("Routes", func(t *testing.T) {
		input := newRouter()
		input.handlers["test"] = HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("test")
		})

		assert.Equal(t, input.handlers, input.Routes())
	})

	t.Run("Middlewares", func(t *testing.T) {
		input := newRouter()
		input.middleware = append(input.middleware, MiddlewareFunction(func(next Handler) Handler { return next }))

		assert.Equal(t, input.middleware, input.Middlewares())
	})

	t.Run("Use", func(t *testing.T) {
		router := newRouter()

		middleware := MiddlewareFunction(func(next Handler) Handler {
			return next
		})

		assert.Len(t, router.middleware, 0)
		router.Use(middleware)
		assert.Len(t, router.middleware, 1)
	})

	t.Run("Help", func(t *testing.T) {
		router := newRouter()
		router.Help(HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("help")
		}))

		assert.NotNil(t, router.helpHander)
		actual := router.helpHander.Execute(nil, nil)
		assert.Equal(t, fmt.Errorf("help"), actual)
	})

	t.Run("NotFound", func(t *testing.T) {
		router := newRouter()
		router.NotFound(HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("not found")
		}))

		assert.NotNil(t, router.notFoundHandler)
		actual := router.notFoundHandler.Execute(nil, nil)
		assert.Equal(t, fmt.Errorf("not found"), actual)
	})
}

func Test_Router_Execute(t *testing.T) {
	// TODO : add tests for flag logic

	router := newRouter()

	t.Run("empty", func(t *testing.T) {
		request := newRequest(context.Background(), []string{"anything"}, nil)
		actual := router.Execute(nil, request)
		assert.Nil(t, actual)
	})

	router.notFoundHandler = HandlerFunction(func(ResponseWriter, Request) error {
		return fmt.Errorf("not found")
	})
	router.handlers["test"] = HandlerFunction(func(ResponseWriter, Request) error {
		return fmt.Errorf("found")
	})

	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "found",
			input:    "test",
			expected: fmt.Errorf("found"),
		},
		{
			name:     "not found",
			input:    "other",
			expected: fmt.Errorf("not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := newRequest(context.Background(), []string{test.input}, nil)
			actual := router.Execute(nil, request)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_Router_Match(t *testing.T) {

	router := newRouter()
	router.Use(MiddlewareFunction(func(next Handler) Handler {
		return HandlerFunction(func(rw ResponseWriter, r Request) error {
			err := next.Execute(rw, r)
			return fmt.Errorf("%w with top middleware", err)
		})
	}))
	router.HandleFunction("found", func(ResponseWriter, Request) error {
		return fmt.Errorf("found")
	})
	router.Route("child", func(r Router) {
		r.HandleFunction("func", func(ResponseWriter, Request) error {
			return fmt.Errorf("found in child")
		})
		r.NotFound(HandlerFunction(func(ResponseWriter, Request) error {
			return fmt.Errorf("not found")
		}))
	})
	router.Group(func(r Router) {
		r.Use(MiddlewareFunction(func(next Handler) Handler {
			return HandlerFunction(func(rw ResponseWriter, r Request) error {
				err := next.Execute(rw, r)
				return fmt.Errorf("%w with group middleware", err)
			})
		}))
		r.HandleFunction("group", func(ResponseWriter, Request) error {
			return fmt.Errorf("found in group")
		})
	})

	type expected struct {
		err   error
		found bool
	}

	tests := []struct {
		name     string
		input    []string
		expected expected
	}{
		{
			name:  "nil",
			input: nil,
			expected: expected{
				err:   nil,
				found: false,
			},
		},
		{
			name:  "none",
			input: []string{},
			expected: expected{
				err:   nil,
				found: false,
			},
		},
		{
			name:  "found",
			input: []string{"found"},
			expected: expected{
				err:   fmt.Errorf("found with top middleware"),
				found: true,
			},
		},
		{
			name:  "child",
			input: []string{"child"},
			expected: expected{
				err:   fmt.Errorf("not found with top middleware"),
				found: true,
			},
		},
		{
			name:  "child func",
			input: []string{"child", "func"},
			expected: expected{
				err:   fmt.Errorf("found in child with top middleware"),
				found: true,
			},
		},
		{
			name:  "func",
			input: []string{"func"},
			expected: expected{
				err:   nil,
				found: false,
			},
		},
		{
			name:  "group",
			input: []string{"group"},
			expected: expected{
				err:   fmt.Errorf("found in group with group middleware with top middleware"),
				found: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler, found := router.Match(test.input)
			assert.Equal(t, test.expected.found, found)

			if handler != nil {
				request := newRequest(context.Background(), test.input[1:], router)
				err := handler.Execute(nil, request)

				if test.expected.err != nil {
					assert.EqualValues(t, test.expected.err.Error(), err.Error())
				} else {
					assert.Nil(t, err)
				}

			} else {
				assert.Nil(t, test.expected.err)
				assert.False(t, found)
			}
		})
	}
}

func Test_Router_Group(t *testing.T) {
	router := newRouter()
	actual := router.Group(func(r Router) {
		r.HandleFunction("found", func(ResponseWriter, Request) error {
			return fmt.Errorf("found")
		})
	})

	assert.Contains(t, router.children, actual)

	request := newRequest(context.Background(), []string{"found"}, nil)
	direct := actual.Execute(nil, request)
	parent := router.Execute(nil, request)

	assert.Equal(t, direct, parent)
	assert.Equal(t, direct, fmt.Errorf("found"))
}

func Test_Router_Route(t *testing.T) {
	t.Run("set route", func(t *testing.T) {
		router := newRouter()
		subRouter := router.Route("route", func(r Router) {
			r.HandleFunction("found", func(ResponseWriter, Request) error {
				return fmt.Errorf("found")
			})
		})

		assert.Contains(t, router.handlers, "route")

		request := newRequest(context.Background(), []string{"route", "found"}, nil)
		parent := router.Execute(nil, request)

		direct := subRouter.Execute(nil, request.WithArgs([]string{"found"}))

		assert.Equal(t, direct, parent)
		assert.Equal(t, direct, fmt.Errorf("found"))
	})

	t.Run("duplicate panic", func(t *testing.T) {
		testPanic(t, func() {
			router := newRouter()
			router.Route("route", func(r Router) {
				r.HandleFunction("first", func(ResponseWriter, Request) error {
					return fmt.Errorf("found")
				})
			})
			router.Route("route", func(r Router) {
				r.HandleFunction("second", func(ResponseWriter, Request) error {
					return fmt.Errorf("found")
				})
			})
		}, errors.DuplicateCommand("route").Error())
	})
}

func Test_Router_Handle(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		router := newRouter()
		router.Handle("found", &testHandler{})
		assert.Contains(t, router.handlers, "found")
	})
	t.Run("duplicate", func(t *testing.T) {
		testPanic(t, func() {
			router := newRouter()
			router.Handle("found", &testHandler{})
			router.Handle("found", &testHandler{})
		}, errors.DuplicateCommand("found").Error())
	})
}

func Test_Router_HandleFunction(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		router := newRouter()
		router.HandleFunction("found", HandlerFunction(func(ResponseWriter, Request) error {
			return nil
		}))
		assert.Contains(t, router.handlers, "found")
	})
	t.Run("duplicate", func(t *testing.T) {
		testPanic(t, func() {
			router := newRouter()
			router.HandleFunction("found", HandlerFunction(func(ResponseWriter, Request) error {
				return nil
			}))
			router.HandleFunction("found", HandlerFunction(func(ResponseWriter, Request) error {
				return nil
			}))
		}, errors.DuplicateCommand("found").Error())
	})
}

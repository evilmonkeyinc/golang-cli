package shell

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/errors"
	"github.com/evilmonkeyinc/golang-cli/flags"
	"github.com/stretchr/testify/assert"
)

// testPanic is a helper function so we can unit test functions that panic
func testPanic(t *testing.T, fn func(), expectedErrMsg string) {
	defer func() {
		if rvr := recover(); rvr != nil {
			errMsg := fmt.Sprintf("%v", rvr)
			switch v := rvr.(type) {
			case error:
				errMsg = v.Error()
			case string:
				errMsg = v
			default:
				break
			}
			assert.Equal(t, expectedErrMsg, errMsg)
			return
		}
		assert.Fail(t, "function should have panicked")
	}()
	fn()
}

func Test_Shell(t *testing.T) {

	t.Run("default setup", func(t *testing.T) {
		actual := &Shell{}
		actual.setup()

		assert.NotNil(t, actual.router)
		assert.NotNil(t, actual.outputWriter)
		assert.Equal(t, os.Stdout, actual.outputWriter)
		assert.NotNil(t, actual.errorWriter)
		assert.Equal(t, os.Stderr, actual.errorWriter)
		assert.NotNil(t, actual.reader)
		assert.Equal(t, os.Stdin, actual.reader)
		assert.NotNil(t, actual.closed)
	})

	t.Run("modified setup", func(t *testing.T) {

		testWritter := &bytes.Buffer{}
		testReader := strings.NewReader("test request")
		testRouter := newRouter()

		actual := &Shell{}
		actual.router = testRouter
		actual.errorWriter = testWritter
		actual.outputWriter = testWritter
		actual.reader = testReader
		actual.setup()

		assert.NotNil(t, actual.router)
		assert.Equal(t, testRouter, actual.router)
		assert.NotNil(t, actual.outputWriter)
		assert.Equal(t, testWritter, actual.outputWriter)
		assert.NotNil(t, actual.errorWriter)
		assert.Equal(t, testWritter, actual.errorWriter)
		assert.NotNil(t, actual.reader)
		assert.Equal(t, testReader, actual.reader)
		assert.NotNil(t, actual.closed)
	})

	t.Run("NotFound", func(t *testing.T) {

		testRouter := newRouter()

		actual := &Shell{}
		actual.router = testRouter

		assert.Nil(t, testRouter.notFoundHandler)
		actual.NotFound(HandlerFunction(func(ResponseWriter, *Request) error {
			return fmt.Errorf("not found")
		}))
		assert.NotNil(t, testRouter.notFoundHandler)
	})

	t.Run("Flags", func(t *testing.T) {

		testRouter := newRouter()

		actual := &Shell{}
		actual.router = testRouter

		assert.Nil(t, testRouter.flags)
		actual.Flags(flags.FlagHandlerFunction(func(fd flags.FlagDefiner) {

		}))
		assert.NotNil(t, testRouter.flags)
	})

	t.Run("Missing Flags", func(t *testing.T) {

		testRouter := newRouter()
		errWriter := &bytes.Buffer{}

		actual := &Shell{}
		actual.router = testRouter
		actual.errorWriter = errWriter

		assert.Nil(t, testRouter.flags)
		actual.Flags(flags.FlagHandlerFunction(func(fd flags.FlagDefiner) {
			fd.Bool("found", false, "")
		}))
		assert.NotNil(t, testRouter.flags)

		actual.execute(context.Background(), []string{"-found", "-missing"})
		assert.Equal(t, "flagset parse failed flag provided but not defined: -missing\n", errWriter.String())
	})

}

func Test_Shell_execute(t *testing.T) {
	type input struct {
		args []string
	}

	type expected struct {
		err error
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "not found",
			input: input{
				args: []string{"invalid"},
			},
			expected: expected{
				err: fmt.Errorf("command not found"),
			},
		},
		{
			name: "found",
			input: input{
				args: []string{"test"},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			shell := &Shell{}
			shell.HandleFunction("test", func(rw ResponseWriter, r *Request) error {
				return nil
			})
			shell.NotFound(HandlerFunction(func(rw ResponseWriter, r *Request) error {
				return fmt.Errorf("command not found")
			}))

			err := shell.execute(context.Background(), test.input.args)

			if test.expected.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, test.expected.err, err)
			}
		})
	}
}

func Test_Shell_Options(t *testing.T) {

	panicOption := OptionFunction(func(shell *Shell) error {
		panic("this function should never be called")
	})

	counter := 0
	validOption := OptionFunction(func(shell *Shell) error {
		counter++
		return nil
	})

	type expected struct {
		err   error
		count int
	}

	tests := []struct {
		name     string
		input    []Option
		expected expected
	}{
		{
			name: "fail first",
			input: []Option{
				OptionFunction(func(shell *Shell) error {
					return fmt.Errorf("fail")
				}),
				panicOption,
			},
			expected: expected{
				err: fmt.Errorf("fail"),
			},
		},
		{
			name: "fail second",
			input: []Option{
				validOption,
				OptionFunction(func(shell *Shell) error {
					return fmt.Errorf("fail")
				}),
				panicOption,
			},
			expected: expected{
				count: 1,
				err:   fmt.Errorf("fail"),
			},
		},
		{
			name: "successful",
			input: []Option{
				validOption,
			},
			expected: expected{
				count: 1,
			},
		},
		{
			name: "three successes",
			input: []Option{
				validOption,
				validOption,
				validOption,
			},
			expected: expected{
				count: 3,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counter = 0
			shell := &Shell{}
			actualErr := shell.Options(test.input...)

			if test.expected.err == nil {
				assert.Nil(t, actualErr)
			} else {
				assert.NotNil(t, actualErr)
				assert.Equal(t, test.expected.err, actualErr)
			}
			assert.Equal(t, test.expected.count, counter)
		})
	}
}

func Test_Shell_Use(t *testing.T) {

	type contextKey string
	var valuesKey contextKey = "values"

	sampleMiddleware := MiddlewareFunction(func(next Handler) Handler {
		return HandlerFunction(func(rw ResponseWriter, r *Request) error {

			ctx := r.Context()

			values, ok := r.Context().Value(valuesKey).([]string)
			assert.True(t, ok, "values should cast to string array")
			ctx = context.WithValue(ctx, valuesKey, append(values, fmt.Sprintf("%d", len(values)+1)))

			r = r.WithContext(ctx)
			return next.Execute(rw, r)
		})
	})

	type expected struct {
		values []string
	}

	tests := []struct {
		name     string
		input    []Middleware
		expected expected
	}{
		{
			name:  "none",
			input: nil,
			expected: expected{
				values: []string{},
			},
		},
		{
			name: "single",
			input: []Middleware{
				sampleMiddleware,
			},
			expected: expected{
				values: []string{
					"1",
				},
			},
		},
		{
			name: "multiple",
			input: []Middleware{
				sampleMiddleware,
				sampleMiddleware,
				sampleMiddleware,
			},
			expected: expected{
				values: []string{
					"1",
					"2",
					"3",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shell := &Shell{}
			shell.Use(test.input...)

			shell.HandleFunction("test", func(rw ResponseWriter, r *Request) error {
				values, ok := r.Context().Value(valuesKey).([]string)
				assert.True(t, ok, "values should cast to string array")
				assert.Equal(t, test.expected.values, values)
				return nil
			})

			ctx := context.WithValue(context.Background(), valuesKey, []string{})
			shell.execute(ctx, []string{"test"})
		})
	}
}

func Test_Shell_Group(t *testing.T) {

	successCommand := "success"
	successHandler := HandlerFunction(func(rw ResponseWriter, r *Request) error {
		return fmt.Errorf("success")
	})

	failureHandler := HandlerFunction(func(rw ResponseWriter, r *Request) error {
		return fmt.Errorf("failure")
	})

	tests := []struct {
		name     string
		input    []func(r Router)
		expected error
	}{
		{
			name: "first group",
			input: []func(Router){
				func(r Router) {
					r.HandleFunction(successCommand, successHandler)
				},
				func(r Router) {
					r.HandleFunction("fail", failureHandler)
				},
			},
			expected: fmt.Errorf("success"),
		},
		{
			name: "second group",
			input: []func(Router){
				func(r Router) {
					r.HandleFunction("fail", failureHandler)
				},
				func(r Router) {
					r.HandleFunction(successCommand, successHandler)
				},
			},
			expected: fmt.Errorf("success"),
		},
		{
			name: "not found",
			input: []func(Router){
				func(r Router) {
					r.HandleFunction("fail", failureHandler)
				},
				func(r Router) {
					r.HandleFunction("other", successHandler)
				},
			},
			expected: fmt.Errorf("not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shell := &Shell{}
			shell.NotFound(HandlerFunction(func(rw ResponseWriter, r *Request) error {
				return fmt.Errorf("not found")
			}))
			for _, group := range test.input {
				shell.Group(group)
			}
			actual := shell.execute(context.Background(), []string{successCommand})
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_Shell_Route(t *testing.T) {

	validArgs := []string{"valid", "test"}

	type route struct {
		command string
		router  func(r Router)
	}

	tests := []struct {
		name     string
		input    []route
		expected error
	}{
		{
			name: "valid",
			input: []route{
				{
					command: "valid",
					router: func(r Router) {
						r.HandleFunction("test", func(rw ResponseWriter, r *Request) error {
							return fmt.Errorf("expected")
						})
					},
				},
			},
			expected: fmt.Errorf("expected"),
		},
		{
			name: "match route but missing command",
			input: []route{
				{
					command: "valid",
					router: func(r Router) {
						r.HandleFunction("missing", func(rw ResponseWriter, r *Request) error {
							return fmt.Errorf("expected")
						})
					},
				},
			},
			expected: fmt.Errorf("not found"),
		},
		{
			name: "match route but missing command custom not found",
			input: []route{
				{
					command: "valid",
					router: func(r Router) {
						r.HandleFunction("missing", func(rw ResponseWriter, r *Request) error {
							return fmt.Errorf("expected")
						})
						r.NotFound(HandlerFunction(func(rw ResponseWriter, r *Request) error {
							return fmt.Errorf("custom not found")
						}))
					},
				},
			},
			expected: fmt.Errorf("custom not found"),
		},
		{
			name: "missing route",
			input: []route{
				{
					command: "invalid",
					router: func(r Router) {
						r.HandleFunction("test", func(rw ResponseWriter, r *Request) error {
							return fmt.Errorf("expected")
						})
					},
				},
			},
			expected: fmt.Errorf("not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shell := &Shell{}
			shell.NotFound(HandlerFunction(func(rw ResponseWriter, r *Request) error {
				return fmt.Errorf("not found")
			}))
			for _, route := range test.input {
				shell.Route(route.command, route.router)
			}
			actual := shell.execute(context.Background(), validArgs)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_Shell_Handle(t *testing.T) {

	type handle struct {
		command string
		handler Handler
	}

	tests := []struct {
		name     string
		input    []handle
		expected error
	}{
		{
			name: "found",
			input: []handle{
				{
					command: "test",
					handler: &testHandler{"found"},
				},
			},
			expected: fmt.Errorf("found"),
		},
		{
			name: "not found",
			input: []handle{
				{
					command: "missing",
					handler: &testHandler{"found"},
				},
			},
			expected: fmt.Errorf("not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shell := &Shell{}
			shell.NotFound(HandlerFunction(func(rw ResponseWriter, r *Request) error {
				return fmt.Errorf("not found")
			}))
			for _, handle := range test.input {
				shell.Handle(handle.command, handle.handler)
			}
			actual := shell.execute(context.Background(), []string{"test"})
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_Shell_HandleFunction(t *testing.T) {

	type handle struct {
		command string
		handler HandlerFunction
	}

	tests := []struct {
		name     string
		input    []handle
		expected error
	}{
		{
			name: "found",
			input: []handle{
				{
					command: "test",
					handler: func(rw ResponseWriter, r *Request) error {
						return fmt.Errorf("found")
					},
				},
			},
			expected: fmt.Errorf("found"),
		},
		{
			name: "not found",
			input: []handle{
				{
					command: "missing",
					handler: func(rw ResponseWriter, r *Request) error {
						return fmt.Errorf("found")
					},
				},
			},
			expected: fmt.Errorf("not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			shell := &Shell{}
			shell.NotFound(HandlerFunction(func(rw ResponseWriter, r *Request) error {
				return fmt.Errorf("not found")
			}))
			for _, handle := range test.input {
				shell.HandleFunction(handle.command, handle.handler)
			}
			actual := shell.execute(context.Background(), []string{"test"})
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_Shell_Execute(t *testing.T) {
	type input struct {
		args []string
	}

	type expected struct {
		args []string
		err  error
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "not found",
			input: input{
				args: []string{"cmd", "invalid"},
			},
			expected: expected{
				args: []string{"invalid"},
				err:  fmt.Errorf("command not found"),
			},
		},
		{
			name: "found",
			input: input{
				args: []string{"cmd", "test"},
			},
			expected: expected{
				args: []string{},
				err:  nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			shell := &Shell{}
			shell.HandleFunction("test", func(rw ResponseWriter, r *Request) error {
				args := r.Args
				assert.Equal(t, test.expected.args, args)
				return nil
			})
			shell.NotFound(HandlerFunction(func(rw ResponseWriter, r *Request) error {
				args := r.Args
				assert.Equal(t, test.expected.args, args)
				return fmt.Errorf("command not found")
			}))

			os.Args = test.input.args
			err := shell.Execute(context.Background())

			if test.expected.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, test.expected.err, err)
			}
		})
	}
}

func Test_Shell_Start(t *testing.T) {
	type input struct {
		args []string
	}

	type expected struct {
		args []string
		err  error
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "not found",
			input: input{
				args: []string{"invalid"},
			},
			expected: expected{
				args: []string{"invalid"},
				err:  fmt.Errorf("command not found"),
			},
		},
		{
			name: "found",
			input: input{
				args: []string{"test", "exit"},
			},
			expected: expected{
				args: []string{},
				err:  nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctx, cancel := context.WithCancel(context.Background())

			testReader := strings.NewReader(strings.Join(test.input.args, "\n") + "\n")
			testOutputWriter := &bytes.Buffer{}
			testErrorWriter := &bytes.Buffer{}

			shell := &Shell{
				reader:       testReader,
				outputWriter: testOutputWriter,
				errorWriter:  testErrorWriter,
				exitOnError:  true,
			}
			shell.HandleFunction("exit", func(rw ResponseWriter, r *Request) error {
				args := r.Args
				assert.Equal(t, test.expected.args, args)
				cancel()
				return nil
			})
			shell.HandleFunction("test", func(rw ResponseWriter, r *Request) error {
				args := r.Args
				assert.Equal(t, test.expected.args, args)
				return nil
			})
			shell.NotFound(HandlerFunction(func(rw ResponseWriter, r *Request) error {
				args := r.Args
				assert.Equal(t, test.expected.args, args)
				return fmt.Errorf("command not found")
			}))

			go func() {
				err := shell.Start(ctx)
				if test.expected.err == nil {
					assert.Nil(t, err)
				} else {
					assert.NotNil(t, err)
					assert.Equal(t, test.expected.err, err)
				}
			}()

			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			select {
			case <-signals:
				cancel()
				<-shell.Closed()
			case <-shell.Closed():
				cancel()
			case <-ctx.Done():
				<-shell.Closed()
			}

		})
	}
}

func Test_Shell_ExitOnError(t *testing.T) {

	t.Run("default", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		testReader := strings.NewReader("error\nexit\n")
		testOutputWriter := &bytes.Buffer{}
		testErrorWriter := &bytes.Buffer{}

		shell := &Shell{
			reader:       testReader,
			errorWriter:  testErrorWriter,
			outputWriter: testOutputWriter,
		}
		shell.HandleFunction("error", func(rw ResponseWriter, r *Request) error {
			return fmt.Errorf("error response")
		})
		shell.HandleFunction("exit", func(rw ResponseWriter, r *Request) error {
			cancel()
			return nil
		})

		go func() {
			os.Args = []string{}
			err := shell.Start(ctx)
			assert.Nil(t, err)
		}()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-signals:
			cancel()
			<-shell.Closed()
		case <-shell.Closed():
			cancel()
		}
	})

	t.Run("false", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		testReader := strings.NewReader("error\nexit\n")
		testOutputWriter := &bytes.Buffer{}
		testErrorWriter := &bytes.Buffer{}

		shell := &Shell{
			reader:       testReader,
			errorWriter:  testErrorWriter,
			outputWriter: testOutputWriter,
			exitOnError:  false,
		}
		shell.HandleFunction("error", func(rw ResponseWriter, r *Request) error {
			return fmt.Errorf("error response")
		})
		shell.HandleFunction("exit", func(rw ResponseWriter, r *Request) error {
			cancel()
			return nil
		})

		go func() {
			os.Args = []string{}
			err := shell.Start(ctx)
			assert.Nil(t, err)
		}()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-signals:
			cancel()
			<-shell.Closed()
		case <-shell.Closed():
			cancel()
		}
	})

	t.Run("true", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		testReader := strings.NewReader("error\nexit\n")
		testOutputWriter := &bytes.Buffer{}
		testErrorWriter := &bytes.Buffer{}

		shell := &Shell{
			reader:       testReader,
			errorWriter:  testErrorWriter,
			outputWriter: testOutputWriter,
			exitOnError:  true,
		}
		shell.HandleFunction("error", func(rw ResponseWriter, r *Request) error {
			return fmt.Errorf("error response")
		})
		shell.HandleFunction("exit", func(rw ResponseWriter, r *Request) error {
			cancel()
			return nil
		})

		go func() {
			os.Args = []string{}
			err := shell.Start(ctx)
			assert.NotNil(t, err)
			assert.Error(t, err, "error response")
		}()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-signals:
			cancel()
			<-shell.Closed()
		case <-shell.Closed():
			cancel()
		}
	})
}

func Test_Shell_HelpFallback(t *testing.T) {

	shell := &Shell{}
	shell.Flags(flags.FlagHandlerFunction(func(fd flags.FlagDefiner) {
		fd.Bool("toUpper", false, "")
	}))
	shell.HandleFunction("help", func(ResponseWriter, *Request) error {
		return errors.HelpRequested("")
	})

	shell.helpHandler = HandlerFunction(func(ResponseWriter, *Request) error {
		return fmt.Errorf("help handler was called")
	})

	t.Run("short help flag", func(t *testing.T) {
		actual := shell.execute(context.Background(), []string{"-h", "ping"})
		assert.EqualError(t, actual, "help handler was called")
	})

	t.Run("long help flag", func(t *testing.T) {
		actual := shell.execute(context.Background(), []string{"-help", "ping"})
		assert.EqualError(t, actual, "help handler was called")
	})

	t.Run("handler returning help", func(t *testing.T) {
		actual := shell.execute(context.Background(), []string{"help", "ping"})
		assert.EqualError(t, actual, "help handler was called")
	})
}

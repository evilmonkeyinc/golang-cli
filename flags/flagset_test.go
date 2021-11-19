package flags

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Validate the DefaultFlagSet struct matches the FlagSet interface
var _ FlagSet = &DefaultFlagSet{}

// Validate the DefaultFlagSet struct matches the FlagDefiner interface
var _ FlagDefiner = &DefaultFlagSet{}

// Validate the DefaultFlagSet struct matches the FlagValues interface
var _ FlagValues = &DefaultFlagSet{}

func Test_NewDefaultFlagSet(t *testing.T) {
	actual := NewDefaultFlagSet()
	assert.IsType(t, &DefaultFlagSet{}, actual)
	assert.Nil(t, actual.set)
}

func Test_NewDefaultFlagSetWithBase(t *testing.T) {
	flagSet := flag.NewFlagSet("", 0)
	actual := NewDefaultFlagSetWithBase(flagSet)
	assert.IsType(t, &DefaultFlagSet{}, actual)
	assert.Equal(t, flagSet, actual.set)
}

func Test_DefaultFlagSet_setup(t *testing.T) {

	t.Run("nil", func(t *testing.T) {
		actual := &DefaultFlagSet{}
		assert.Nil(t, actual.set)
		actual.setup()
		assert.NotNil(t, actual.set)
	})

	t.Run("set", func(t *testing.T) {
		flagSet := flag.NewFlagSet("", 0)
		actual := &DefaultFlagSet{
			set: flagSet,
		}
		assert.NotNil(t, actual.set)
		assert.Equal(t, flagSet, actual.set)
		actual.setup()
		assert.NotNil(t, actual.set)
		assert.Equal(t, flagSet, actual.set)
	})
}

func Test_DefaultFlagSet_SubFlagSet(t *testing.T) {
	type input struct {
		args     []string
		handlers []FlagHandlerFunction
	}

	type expected struct {
		args       []string
		flagValues map[string]interface{}
	}

	tests := []struct {
		name     string
		input    input
		expected []expected
	}{
		{
			name: "one",
			input: input{
				args: []string{"-bool", "arg"},
				handlers: []FlagHandlerFunction{
					func(fd FlagDefiner) {
						fd.Bool("bool", false, "a bool")
					},
				},
			},
			expected: []expected{
				{
					args: []string{"arg"},
					flagValues: map[string]interface{}{
						"bool": true,
					},
				},
			},
		},
		{
			name: "two split",
			input: input{
				args: []string{"-bool", "arg", "-string", "value", "arg"},
				handlers: []FlagHandlerFunction{
					func(fd FlagDefiner) {
						fd.Bool("bool", false, "a bool")
					},
					func(fd FlagDefiner) {
						fd.String("string", "", "a string")
					},
				},
			},
			expected: []expected{
				{
					args: []string{"arg", "-string", "value", "arg"},
					flagValues: map[string]interface{}{
						"bool": true,
					},
				},
				{
					args: []string{"arg"},
					flagValues: map[string]interface{}{
						"bool":   true,
						"string": "value",
					},
				},
			},
		},
		{
			name: "two defined first but split given",
			input: input{
				args: []string{"-bool", "arg", "-string", "value", "arg"},
				handlers: []FlagHandlerFunction{
					func(fd FlagDefiner) {
						fd.Bool("bool", false, "a bool")
						fd.String("string", "def", "a string")
					},
					func(fd FlagDefiner) {
					},
				},
			},
			expected: []expected{
				{
					args: []string{"arg", "-string", "value", "arg"},
					flagValues: map[string]interface{}{
						"bool":   true,
						"string": "def",
					},
				},
				{
					args: []string{"arg"},
					flagValues: map[string]interface{}{
						"bool":   true,
						"string": "value",
					},
				},
			},
		},
		{
			name: "two defined first but given after",
			input: input{
				args: []string{"arg", "-bool", "-string", "value", "arg"},
				handlers: []FlagHandlerFunction{
					func(fd FlagDefiner) {
						fd.Bool("bool", false, "a bool")
						fd.String("string", "def", "a string")
					},
					func(fd FlagDefiner) {
					},
				},
			},
			expected: []expected{
				{
					args: []string{"arg", "-bool", "-string", "value", "arg"},
					flagValues: map[string]interface{}{
						"bool":   false,
						"string": "def",
					},
				},
				{
					args: []string{"arg"},
					flagValues: map[string]interface{}{
						"bool":   true,
						"string": "value",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var flagSet FlagSet = &DefaultFlagSet{
				set: flag.NewFlagSet("", flag.ContinueOnError),
			}
			args := test.input.args

			for i, fn := range test.input.handlers {
				fn(flagSet)
				args, _ = flagSet.Parse(args)

				expected := test.expected[i]
				assert.Equal(t, expected.args, args)

				for key, value := range expected.flagValues {
					actual := flagSet.Get(key)
					assert.Equal(t, value, actual)
				}

				// advance command line
				assert.Equal(t, "arg", args[0])
				args = args[1:]
				flagSet = flagSet.SubFlagSet(fmt.Sprintf("%d", i))
			}
		})
	}
}

func Test_DefaultFlagSet_ParseArgs(t *testing.T) {

	tests := []struct {
		name          string
		input         []string
		expected      []string
		expectedError error
	}{
		{
			name:     "one arg no flags",
			input:    []string{"ping"},
			expected: []string{"ping"},
		},
		{
			name:     "one flag one arg",
			input:    []string{"--one=one", "ping"},
			expected: []string{"ping"},
		},
		{
			name:     "one flag one arg alt",
			input:    []string{"--one", "one", "ping"},
			expected: []string{"ping"},
		},
		{
			name:     "one flag one arg one flag",
			input:    []string{"--one", "one", "ping", "--two", "two"},
			expected: []string{"ping", "--two", "two"},
		},
		{
			name:          "short help",
			input:         []string{"-h"},
			expected:      []string{},
			expectedError: ErrHelp,
		},
		{
			name:          "long help",
			input:         []string{"--help"},
			expected:      []string{},
			expectedError: ErrHelp,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			flagSet := &DefaultFlagSet{}
			flagSet.String("one", "", "")
			flagSet.String("two", "", "")
			flagSet.String("three", "", "")

			assert.False(t, flagSet.Parsed())

			actual, actualError := flagSet.Parse(test.input)

			assert.True(t, flagSet.Parsed())

			assert.Equal(t, test.expected, actual)
			assert.Equal(t, test.expectedError, actualError)
		})
	}
}

func Test_DefaultFlagSet_Set(t *testing.T) {

	type input struct {
		name  string
		value string
		fn    FlagHandlerFunction
	}

	type expected struct {
		err   error
		value interface{}
	}

	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "missing",
			input: input{
				name:  "missing",
				value: "anything",
				fn:    func(fd FlagDefiner) {},
			},
			expected: expected{
				err:   fmt.Errorf("no such flag -missing"),
				value: nil,
			},
		},
		{
			name: "valid bool",
			input: input{
				name:  "bool",
				value: "true",
				fn: func(fd FlagDefiner) {
					fd.Bool("bool", false, "")
				},
			},
			expected: expected{
				err:   nil,
				value: true,
			},
		},
		{
			name: "invalid bool",
			input: input{
				name:  "bool",
				value: "ok",
				fn: func(fd FlagDefiner) {
					fd.Bool("bool", false, "")
				},
			},
			expected: expected{
				err:   fmt.Errorf("parse error"),
				value: false,
			},
		},
		{
			name: "valid string",
			input: input{
				name:  "string",
				value: "string value",
				fn: func(fd FlagDefiner) {
					fd.String("string", "", "")
				},
			},
			expected: expected{
				err:   nil,
				value: "string value",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flagSet := &DefaultFlagSet{}
			test.input.fn(flagSet)

			actualErr := flagSet.Set(test.input.name, test.input.value)
			assert.Equal(t, test.expected.err, actualErr)

			actualVal := flagSet.Get(test.input.name)
			assert.Equal(t, test.expected.value, actualVal)
		})
	}
}

type customFlagValue string

func (v *customFlagValue) String() string {
	return string(*v)
}
func (v *customFlagValue) Set(val string) error {
	*v = customFlagValue(val)
	return nil
}

func Test_DefaultFlagSet_Get(t *testing.T) {

	t.Run("GetBool", func(t *testing.T) {
		flagSet := &DefaultFlagSet{}
		flagSet.Bool("valid", false, "")
		flagSet.String("invalid", "true", "")

		value, ok := flagSet.GetBool("valid")
		assert.True(t, ok)
		assert.Equal(t, false, value)

		flagSet.Set("valid", "true")
		value, ok = flagSet.GetBool("valid")
		assert.True(t, ok)
		assert.Equal(t, true, value)

		value, ok = flagSet.GetBool("invalid")
		assert.False(t, ok)
		assert.Equal(t, false, value)

		value, ok = flagSet.GetBool("missing")
		assert.False(t, ok)
		assert.Equal(t, false, value)
	})

	t.Run("GetInt", func(t *testing.T) {
		flagSet := &DefaultFlagSet{}
		flagSet.Int("valid", 10, "")
		flagSet.String("invalid", "10", "")

		value, ok := flagSet.GetInt("valid")
		assert.Equal(t, int64(10), value)
		assert.True(t, ok)

		value, ok = flagSet.GetInt("invalid")
		assert.Equal(t, int64(0), value)
		assert.False(t, ok)

		value, ok = flagSet.GetInt("missing")
		assert.Equal(t, int64(0), value)
		assert.False(t, ok)
	})

	t.Run("GetUint", func(t *testing.T) {
		flagSet := &DefaultFlagSet{}
		flagSet.Uint("valid", 10, "")
		flagSet.String("invalid", "10", "")

		value, ok := flagSet.GetUint("valid")
		assert.Equal(t, uint64(10), value)
		assert.True(t, ok)

		value, ok = flagSet.GetUint("invalid")
		assert.Equal(t, uint64(0), value)
		assert.False(t, ok)

		value, ok = flagSet.GetUint("missing")
		assert.Equal(t, uint64(0), value)
		assert.False(t, ok)
	})

	t.Run("GetString", func(t *testing.T) {
		flagSet := &DefaultFlagSet{}
		flagSet.String("valid", "valid", "")
		flagSet.Bool("invalid", false, "")

		value, ok := flagSet.GetString("valid")
		assert.Equal(t, "valid", value)
		assert.True(t, ok)

		value, ok = flagSet.GetString("invalid")
		assert.Equal(t, "", value)
		assert.False(t, ok)

		value, ok = flagSet.GetString("missing")
		assert.Equal(t, "", value)
		assert.False(t, ok)
	})

	t.Run("GetFloat", func(t *testing.T) {
		flagSet := &DefaultFlagSet{}
		flagSet.Float("valid", 42.1, "")
		flagSet.String("invalid", "10", "")

		value, ok := flagSet.GetFloat("valid")
		assert.Equal(t, float64(42.1), value)
		assert.True(t, ok)

		value, ok = flagSet.GetFloat("invalid")
		assert.Equal(t, float64(0), value)
		assert.False(t, ok)

		value, ok = flagSet.GetFloat("missing")
		assert.Equal(t, float64(0), value)
		assert.False(t, ok)
	})

	t.Run("GetDuration", func(t *testing.T) {
		flagSet := &DefaultFlagSet{}
		flagSet.Duration("valid", time.Hour, "")
		flagSet.String("invalid", "10h", "")

		value, ok := flagSet.GetDuration("valid")
		assert.Equal(t, time.Hour, value)
		assert.True(t, ok)

		value, ok = flagSet.GetDuration("invalid")
		assert.Equal(t, time.Duration(0), value)
		assert.False(t, ok)

		value, ok = flagSet.GetDuration("missing")
		assert.Equal(t, time.Duration(0), value)
		assert.False(t, ok)
	})

	t.Run("Var", func(t *testing.T) {
		arrayFlag := &StringArrayFlag{
			"one", "two", "three",
		}

		flagSet := &DefaultFlagSet{}
		flagSet.Var(arrayFlag, "valid", "")

		valid := flagSet.Get("valid")
		assert.Equal(t, arrayFlag.Get(), valid)
	})

	t.Run("GetCustom", func(t *testing.T) {

		val := customFlagValue("value")

		base := flag.NewFlagSet("", flag.ContinueOnError)
		base.Var(&val, "custom", "")
		flagSet := NewDefaultFlagSetWithBase(base)

		actual := flagSet.Get("custom")
		assert.Equal(t, "value", actual)
	})
}

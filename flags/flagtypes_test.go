package flags

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringArrayFlag(t *testing.T) {

	type expected struct {
		array []string
		str   string
	}

	tests := []struct {
		name     string
		input    []string
		expected expected
	}{
		{
			name:  "one",
			input: []string{"-array", "one"},
			expected: expected{
				array: []string{"one"},
				str:   "one",
			},
		},
		{
			name:  "two",
			input: []string{"-array", "one", "-array", "two"},
			expected: expected{
				array: []string{"one", "two"},
				str:   "one,two",
			},
		},
		{
			name:  "three",
			input: []string{"-array", "one", "-array", "two", "-array", "three"},
			expected: expected{
				array: []string{"one", "two", "three"},
				str:   "one,two,three",
			},
		},
		{
			name:  "comma-seperated",
			input: []string{"-array", "one,two,three"},
			expected: expected{
				array: []string{"one", "two", "three"},
				str:   "one,two,three",
			},
		},
		{
			name:  "combo",
			input: []string{"-array", "one", "-array", "two,three"},
			expected: expected{
				array: []string{"one", "two", "three"},
				str:   "one,two,three",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flagSet := flag.NewFlagSet("", 0)

			value := &StringArrayFlag{}
			flagSet.Var(value, "array", "an array")
			flagSet.Parse(test.input)

			assert.EqualValues(t, test.expected.array, value.Get())
			assert.EqualValues(t, test.expected.str, value.String())
		})
	}
}

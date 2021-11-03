package shell

import (
	"os"
	"strings"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// Validate the OptionFunction func matches the Option interface
var _ Option = OptionFunction(func(shell *Shell) error {
	return nil
})

func Test_OptionInput(t *testing.T) {

	t.Run("not set", func(t *testing.T) {
		reader := strings.NewReader("test")

		option := OptionInput(reader)
		shell := &Shell{}
		err := option.Apply(shell)

		assert.Equal(t, reader, shell.reader)
		assert.Nil(t, err)
	})

	t.Run("already set", func(t *testing.T) {
		reader := strings.NewReader("test")

		option := OptionInput(reader)
		shell := &Shell{
			reader: os.Stdin,
		}
		err := option.Apply(shell)

		assert.NotEqual(t, reader, shell.reader)
		assert.NotNil(t, err)

		expectedError := errors.OptionIsSet("Input")
		assert.EqualValues(t, expectedError, err)
	})

}

func Test_OptionOutputWriter(t *testing.T) {

	t.Run("not set", func(t *testing.T) {
		writer := &testWriter{}

		option := OptionOutputWriter(writer)
		shell := &Shell{}
		err := option.Apply(shell)

		assert.Equal(t, writer, shell.outputWriter)
		assert.Nil(t, err)
	})

	t.Run("already set", func(t *testing.T) {
		writer := &testWriter{}

		option := OptionOutputWriter(writer)
		shell := &Shell{
			outputWriter: os.Stdout,
		}
		err := option.Apply(shell)

		assert.NotEqual(t, writer, shell.outputWriter)
		assert.NotNil(t, err)

		expectedError := errors.OptionIsSet("OutputWriter")
		assert.EqualValues(t, expectedError, err)
	})

}
func Test_OptionErrorWriter(t *testing.T) {

	t.Run("not set", func(t *testing.T) {
		writer := &testWriter{}

		option := OptionErrorWriter(writer)
		shell := &Shell{}
		err := option.Apply(shell)

		assert.Equal(t, writer, shell.errorWriter)
		assert.Nil(t, err)
	})

	t.Run("already set", func(t *testing.T) {
		writer := &testWriter{}

		option := OptionErrorWriter(writer)
		shell := &Shell{
			errorWriter: os.Stdout,
		}
		err := option.Apply(shell)

		assert.NotEqual(t, writer, shell.errorWriter)
		assert.NotNil(t, err)

		expectedError := errors.OptionIsSet("ErrorWriter")
		assert.EqualValues(t, expectedError, err)
	})

}

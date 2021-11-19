package shell

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/evilmonkeyinc/golang-cli/errors"
	"github.com/evilmonkeyinc/golang-cli/flags"
	"github.com/stretchr/testify/assert"
)

// Validate the OptionFunction func matches the Option interface
var _ Option = OptionFunction(func(shell *Shell) error {
	return nil
})

func Test_OptionFunction(t *testing.T) {

	fn := OptionFunction(func(shell *Shell) error {
		return fmt.Errorf("option function")
	})

	direct := fn(nil)
	apply := fn.Apply(nil)

	assert.Equal(t, fmt.Errorf("option function"), direct)
	assert.Equal(t, fmt.Errorf("option function"), apply)
	assert.Equal(t, direct, apply)

}

func Test_OptionInput(t *testing.T) {

	t.Run("invalid", func(t *testing.T) {
		expected := errors.OptionIsInvalid("Input")
		testPanic(t, func() {
			OptionInput(nil)
		}, expected.Error())
	})

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

	t.Run("invalid", func(t *testing.T) {
		expected := errors.OptionIsInvalid("OutputWriter")
		testPanic(t, func() {
			OptionOutputWriter(nil)
		}, expected.Error())
	})

	t.Run("not set", func(t *testing.T) {
		writer := &bytes.Buffer{}

		option := OptionOutputWriter(writer)
		shell := &Shell{}
		err := option.Apply(shell)

		assert.Equal(t, writer, shell.outputWriter)
		assert.Nil(t, err)
	})

	t.Run("already set", func(t *testing.T) {
		writer := &bytes.Buffer{}

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

	t.Run("invalid", func(t *testing.T) {
		expected := errors.OptionIsInvalid("ErrorWriter")
		testPanic(t, func() {
			OptionErrorWriter(nil)
		}, expected.Error())
	})

	t.Run("not set", func(t *testing.T) {
		writer := &bytes.Buffer{}

		option := OptionErrorWriter(writer)
		shell := &Shell{}
		err := option.Apply(shell)

		assert.Equal(t, writer, shell.errorWriter)
		assert.Nil(t, err)
	})

	t.Run("already set", func(t *testing.T) {
		writer := &bytes.Buffer{}

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

func Test_OptionShellPrompt(t *testing.T) {

	t.Run("invalid", func(t *testing.T) {
		expected := errors.OptionIsInvalid("ShellPrompt")
		testPanic(t, func() {
			OptionShellPrompt("")
		}, expected.Error())
	})

	t.Run("not set", func(t *testing.T) {
		prompt := "cli:"

		option := OptionShellPrompt(prompt)
		shell := &Shell{}
		err := option.Apply(shell)

		assert.Equal(t, prompt, shell.shellPrompt)
		assert.Nil(t, err)
	})

	t.Run("already set", func(t *testing.T) {
		prompt := "cli:"

		option := OptionShellPrompt(prompt)
		shell := &Shell{
			shellPrompt: defaultShellPrompt,
		}
		err := option.Apply(shell)

		assert.NotEqual(t, prompt, shell.shellPrompt)
		assert.NotNil(t, err)

		expectedError := errors.OptionIsSet("ShellPrompt")
		assert.EqualValues(t, expectedError, err)
	})

}

func Test_OptionFlagSet(t *testing.T) {

	t.Run("invalid", func(t *testing.T) {
		expected := errors.OptionIsInvalid("FlagSet")
		testPanic(t, func() {
			OptionFlagSet(nil)
		}, expected.Error())
	})

	t.Run("not set", func(t *testing.T) {
		flagSet := flags.NewDefaultFlagSetWithBase(nil)

		option := OptionFlagSet(flagSet)
		shell := &Shell{}
		err := option.Apply(shell)

		assert.Equal(t, flagSet, shell.flagSet)
		assert.Nil(t, err)
	})

	t.Run("already set", func(t *testing.T) {
		flagSet := flags.NewDefaultFlagSetWithBase(flag.CommandLine)

		option := OptionFlagSet(flagSet)
		shell := &Shell{
			flagSet: flags.NewDefaultFlagSetWithBase(nil),
		}
		err := option.Apply(shell)

		assert.NotEqual(t, flagSet, shell.flagSet)
		assert.NotNil(t, err)

		expectedError := errors.OptionIsSet("FlagSet")
		assert.EqualValues(t, expectedError, err)
	})

}

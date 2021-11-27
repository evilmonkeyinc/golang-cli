package commands

import (
	"testing"

	"github.com/evilmonkeyinc/golang-cli/shell"
	"github.com/stretchr/testify/assert"
)

// Validate the CommandRouter struct matches the CommandHandler interface
var _ CommandHandler = &CommandRouter{}

// Validate the CommandRouter struct matches the Router interface
var _ shell.Router = &CommandRouter{}

func Test_CommandRouter(t *testing.T) {

	t.Run("Getters", func(t *testing.T) {
		command := &CommandRouter{
			Name:        "Name",
			Summary:     "The command summary",
			Description: "The command description",
			Usage:       "name <arg1>",
		}
		assert.Equal(t, "Name", command.GetName())
		assert.Equal(t, "The command summary", command.GetSummary())
		assert.Equal(t, "The command description", command.GetDescription())
		assert.Equal(t, "name <arg1>", command.GetUsage())
	})

}

func Test_NewCommandRouter(t *testing.T) {

	called := false
	actual := NewCommandRouter("Name", "The Summary", "The description string", "name <arg1>", func(r shell.Router) {
		assert.IsType(t, &CommandRouter{}, r)
		called = true
	})
	assert.IsType(t, &CommandRouter{}, actual)
	assert.True(t, called, "setup function not called")
}

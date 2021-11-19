package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Validate the FlagHandlerFunction func matches the FlagHandler interface
var _ FlagHandler = FlagHandlerFunction(func(FlagDefiner) {})

func Test_FlagHandlerFunction(t *testing.T) {

	handler := FlagHandlerFunction(func(fd FlagDefiner) {
		fd.Bool("bool", true, "")
	})

	fd := &DefaultFlagSet{}
	handler.Define(fd)

	boolVal, _ := fd.GetBool("bool")
	assert.True(t, boolVal)
}

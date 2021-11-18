package shell

import (
	"bytes"
	"flag"
	"time"
)

var (
	// ErrHelp is the error returned if the -help or -h flag is invoked but no such flag is defined.
	ErrHelp error = flag.ErrHelp
)

// A FlagSet represents a set of defined flags
type FlagSet interface {
	FlagDefiner
	FlagValues

	// SubFlagSet creates a new flagset that will be used by sub routers and sub commands.
	SubFlagSet(name string) FlagSet

	// Parse parses flag definitions from the argument list, which should not include the command name, and return remaining, non-flag, arguments.
	// Must be called after all flags in the FlagSet are defined and before flags are accessed by the program.
	// The return value will be ErrHelp if -help was set but not defined.
	Parse(args []string) ([]string, error)
	// Parsed returns true if Parse has been called.
	Parsed() bool
}

// FlagDefiner allows you to define the flags managed by the flag set
type FlagDefiner interface {
	// Bool defines a bool flag with specified name, default value, and usage string.
	Bool(name string, defaultValue bool, usage string)
	// Int defines a int64 flag with specified name, default value, and usage string.
	Int(name string, defaultValue int64, usage string)
	// Uint defines a unit64 flag with specified name, default value, and usage string.
	Uint(name string, defaultValue uint64, usage string)
	// String defines a string flag with specified name, default value, and usage string.
	String(name string, defaultValue string, usage string)
	// Float defines a float64 flag with specified name, default value, and usage string.
	Float(name string, defaultValue float64, usage string)
	// Duration defines a time.Duration flag with specified name, default value, and usage string.
	Duration(name string, defaultValue time.Duration, usage string)
}

//FlagValues allows you to retreive flags
type FlagValues interface {
	// Get returns the value of the named flag.
	Get(name string) interface{}
	// GetBool returns the value of a named flag as a bool.
	GetBool(name string) *bool
	// GetInt returns the value of a named flag as a int64.
	GetInt(name string) *int64
	// GetUint returns the value of a named flag as a uint64.
	GetUint(name string) *uint64
	// GetString returns the value of a named flag as a string.
	GetString(name string) *string
	// GetFloat returns the value of a named flag as a float64.
	GetFloat(name string) *float64
	// GetDuration returns the value of a named flag as a time.Duration.
	GetDuration(name string) *time.Duration
	// Set sets the value of the named flag.
	Set(name, value string) error
}

// NewDefaultFlagSet returns a new DefaultFlagSet.
func NewDefaultFlagSet() *DefaultFlagSet {
	return &DefaultFlagSet{}
}

// NewDefaultFlagSetWithBase returns a new DefaultFlagSet using the specified flag.FlagSet as a base.
func NewDefaultFlagSetWithBase(flagSet *flag.FlagSet) *DefaultFlagSet {
	return &DefaultFlagSet{
		set: flagSet,
	}
}

// DefaultFlagSet is the basic FlagSet implementation using the standard golang flag library
type DefaultFlagSet struct {
	set *flag.FlagSet
}

func (flagSet *DefaultFlagSet) setup() {
	if flagSet.set == nil {
		flagSet.set = flag.NewFlagSet("", flag.ContinueOnError)
		flagSet.set.SetOutput(&bytes.Buffer{})
	}
}

// SubFlagSet creates a new flagset that will be used by sub routers and sub commands.
func (flagSet *DefaultFlagSet) SubFlagSet(name string) FlagSet {
	flagSet.setup()

	newFlagSet := flag.NewFlagSet(name, flagSet.set.ErrorHandling())
	newFlagSet.SetOutput(flagSet.set.Output())
	newFlagSet.Usage = flagSet.set.Usage

	flagSet.set.VisitAll(func(f *flag.Flag) {
		newFlagSet.Var(f.Value, f.Name, f.Usage)
	})

	return &DefaultFlagSet{
		set: newFlagSet,
	}
}

// Parse parses flag definitions from the argument list, which should not include the command name, and return remaining, non-flag, arguments.
// Must be called after all flags in the FlagSet are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help was set but not defined.
func (flagSet *DefaultFlagSet) Parse(args []string) ([]string, error) {
	flagSet.setup()
	result := flagSet.set.Parse(args)
	return flagSet.set.Args(), result
}

// Parsed returns true if Parse has been called.
func (flagSet *DefaultFlagSet) Parsed() bool {
	flagSet.setup()
	return flagSet.set.Parsed()
}

// Set sets the value of the named flag.
func (flagSet *DefaultFlagSet) Set(name, value string) error {
	flagSet.setup()
	return flagSet.set.Set(name, value)
}

// Get returns the value of the named flag.
func (flagSet *DefaultFlagSet) Get(name string) interface{} {
	flagSet.setup()
	lookedUp := flagSet.set.Lookup(name)
	if lookedUp == nil {
		return nil
	}
	if getter, ok := lookedUp.Value.(flag.Getter); ok {
		return getter.Get()
	}
	return lookedUp.Value.String()
}

// Bool defines a bool flag with specified name, default value, and usage string.
func (flagSet *DefaultFlagSet) Bool(name string, defaultValue bool, usage string) {
	flagSet.setup()
	flagSet.set.Bool(name, defaultValue, usage)
}

// Int defines a int64 flag with specified name, default value, and usage string.
func (flagSet *DefaultFlagSet) Int(name string, defaultValue int64, usage string) {
	flagSet.setup()
	flagSet.set.Int64(name, defaultValue, usage)
}

// Uint defines a unit64 flag with specified name, default value, and usage string.
func (flagSet *DefaultFlagSet) Uint(name string, defaultValue uint64, usage string) {
	flagSet.setup()
	flagSet.set.Uint64(name, defaultValue, usage)
}

// String defines a string flag with specified name, default value, and usage string.
func (flagSet *DefaultFlagSet) String(name string, defaultValue string, usage string) {
	flagSet.setup()
	flagSet.set.String(name, defaultValue, usage)
}

// Float defines a float64 flag with specified name, default value, and usage string.
func (flagSet *DefaultFlagSet) Float(name string, defaultValue float64, usage string) {
	flagSet.setup()
	flagSet.set.Float64(name, defaultValue, usage)
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
func (flagSet *DefaultFlagSet) Duration(name string, defaultValue time.Duration, usage string) {
	flagSet.setup()
	flagSet.set.Duration(name, defaultValue, usage)
}

// GetBool returns the value of a named flag as a bool.
func (flagSet *DefaultFlagSet) GetBool(name string) *bool {
	value := flagSet.Get(name)
	if value == nil {
		return nil
	}
	if boolValue, ok := value.(bool); ok {
		return &boolValue
	}
	return nil
}

// GetInt returns the value of a named flag as a int64.
func (flagSet *DefaultFlagSet) GetInt(name string) *int64 {
	value := flagSet.Get(name)
	if value == nil {
		return nil
	}
	if intValue, ok := value.(int64); ok {
		return &intValue
	}
	return nil
}

// GetUint returns the value of a named flag as a uint64.
func (flagSet *DefaultFlagSet) GetUint(name string) *uint64 {
	value := flagSet.Get(name)
	if value == nil {
		return nil
	}
	if intValue, ok := value.(uint64); ok {
		return &intValue
	}
	return nil
}

// GetString returns the value of a named flag as a string.
func (flagSet *DefaultFlagSet) GetString(name string) *string {
	value := flagSet.Get(name)
	if value == nil {
		return nil
	}
	if stringValue, ok := value.(string); ok {
		return &stringValue
	}
	return nil
}

// GetFloat returns the value of a named flag as a float64.
func (flagSet *DefaultFlagSet) GetFloat(name string) *float64 {
	value := flagSet.Get(name)
	if value == nil {
		return nil
	}
	if floatValue, ok := value.(float64); ok {
		return &floatValue
	}
	return nil
}

// GetDuration returns the value of a named flag as a time.Duration.
func (flagSet *DefaultFlagSet) GetDuration(name string) *time.Duration {
	value := flagSet.Get(name)
	if value == nil {
		return nil
	}
	if durationValue, ok := value.(time.Duration); ok {
		return &durationValue
	}
	return nil
}

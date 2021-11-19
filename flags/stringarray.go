package flags

import "strings"

// StringArrayFlag is a Value implementation that
// adds string array support to the FlagSet
// StringArrayFlag supports multiple uses of the same flag,
// as well as comma seperated lists.
type StringArrayFlag []string

// String retruns the string array as a comma seperated string.
func (flag *StringArrayFlag) String() string {
	return strings.Join(*flag, ",")
}

// Set adds the supplied string to the string array.
// If the string contains comma seperators, the string will be split
// before being added to the array.
func (flag *StringArrayFlag) Set(value string) error {
	split := strings.Split(value, ",")
	*flag = append(*flag, split...)
	return nil
}

// Get returns the string array value.
func (flag *StringArrayFlag) Get() interface{} {
	return []string(*flag)
}

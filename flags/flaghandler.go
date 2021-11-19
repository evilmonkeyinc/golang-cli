package flags

// The FlagHandlerFunction type is an adapter to allow the use of ordinary functions as shell flag handlers.
type FlagHandlerFunction func(FlagDefiner)

// Define allows the function to define command-line
func (handler FlagHandlerFunction) Define(flagDefiner FlagDefiner) {
	handler(flagDefiner)
}

// FlagHandler allows shell handlers to define additional
type FlagHandler interface {
	// Define allows the function to define command-line
	Define(FlagDefiner)
}

package lightmigrate

// Logger is an interface, so you can pass in your own
// logging implementation.
type Logger interface {
	// Printf is like fmt.Printf
	Printf(format string, v ...interface{})
}

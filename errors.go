package lightmigrate

import (
	"fmt"
	"os"
)

type ErrDuplicateMigration struct {
	migration
	os.FileInfo
}

// Error implements error interface.
func (e ErrDuplicateMigration) Error() string {
	return "duplicate migration file: " + e.Name()
}

var (
	ErrDatabaseDirty     = fmt.Errorf("database contains unsucessful migration")
	ErrNoChange          = fmt.Errorf("no change")
	ErrVersionNotAllowed = fmt.Errorf("version 0 is not allowed")
)

// DriverError should be used for errors involving queries ran against the database
type DriverError struct {
	// Optional: the line number
	Line uint

	// Query is a query excerpt
	Query []byte

	// Msg is a useful/helping error message for humans
	Msg string

	// OrigErr is the underlying error
	OrigErr error
}

func (e DriverError) Error() string {
	if e.Msg == "" {
		return fmt.Sprintf("%v in line %v: %s", e.OrigErr, e.Line, e.Query)
	}
	return fmt.Sprintf("%v in line %v: %s (details: %v)", e.Msg, e.Line, e.Query, e.OrigErr)
}

func (e DriverError) Unwrap() error {
	return e.OrigErr
}

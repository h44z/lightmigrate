package lightmigrate

import (
	"io"
)

// MigrationSource is the interface every migration source must implement.
type MigrationSource interface {
	// Closer will clean up the migration source instance.
	io.Closer

	// First returns the very first migration version available.
	// If there is no version available, it must return os.ErrNotExist.
	First() (version uint64, err error)

	// Prev returns the previous version for a given version.
	// If there is no previous version available, it must return os.ErrNotExist.
	Prev(version uint64) (prevVersion uint64, err error)

	// Next returns the next version for a given version.
	// If there is no next version available, it must return os.ErrNotExist.
	Next(version uint64) (nextVersion uint64, err error)

	// ReadUp returns the UP migration body and an identifier that helps
	// finding this migration in the source for a given version.
	// If there is no up migration available for this version,
	// it must return os.ErrNotExist.
	ReadUp(version uint64) (r io.ReadCloser, identifier string, err error)

	// ReadDown returns the DOWN migration body and an identifier that helps
	// finding this migration in the source for a given version.
	// If there is no down migration available for this version,
	// it must return os.ErrNotExist.
	ReadDown(version uint64) (r io.ReadCloser, identifier string, err error)
}

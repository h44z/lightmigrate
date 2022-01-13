package lightmigrate

import "io"

// MigrationDriver is the interface every database driver must implement.
type MigrationDriver interface {
	io.Closer

	// Lock should acquire a database lock so that only one migration process
	// can run at a time. Migrate will call this function before Run is called.
	// If the implementation can't provide this functionality, return nil.
	// Return database.ErrLocked if database is already locked.
	Lock() error

	// Unlock should release the lock. Migrate will call this function after
	// all migrations have been run.
	Unlock() error

	// GetVersion returns the currently active version and the database dirty state.
	// When no migration has been applied, it must return version NoMigrationVersion (0).
	// Dirty means, a previous migration failed and user interaction is required.
	GetVersion() (version uint64, dirty bool, err error)

	// SetVersion saves version and dirty state.
	// Migrate will call this function before and after each call to RunMigration.
	// version must be >= 1. 0 means NoMigrationVersion.
	SetVersion(version uint64, dirty bool) error

	// RunMigration applies a migration to the database. migration is guaranteed to be not nil.
	RunMigration(migration io.Reader) error

	// Reset deletes everything related to LightMigrate in the database.
	Reset() error
}

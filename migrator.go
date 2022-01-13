package lightmigrate

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var NoMigrationVersion uint64 = 0

// Migrator is a generic interface that provides compatibility with golang-migrate migrator.
type Migrator interface {
	Migrate(version uint64) error
}

// migrator contains the main logic for applying migrations.
type migrator struct {
	source   MigrationSource
	driver   MigrationDriver
	lock     sync.Mutex
	shutdown chan bool
	logger   Logger
	verbose  bool
}

type MigratorOption func(svc *migrator)

type migrationData struct {
	Version       uint64
	TargetVersion uint64
	Identifier    string
	Direction     Direction
	Contents      io.ReadCloser

	error error
}

func (m migrationData) Error() error {
	return m.error
}

func NewMigrator(source MigrationSource, driver MigrationDriver, opts ...MigratorOption) (*migrator, error) {
	m := &migrator{
		source: source,
		driver: driver,
		lock:   sync.Mutex{},
		logger: log.Default(),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m, nil
}

func WithLogger(logger Logger) MigratorOption {
	return func(m *migrator) {
		m.logger = logger
	}
}

func WithVerboseLogging(verbose bool) MigratorOption {
	return func(m *migrator) {
		m.verbose = verbose
	}
}

func (m *migrator) Migrate(version uint64) error {
	// avoid multiple concurrent runs of the migration
	m.lock.Lock()
	defer m.lock.Unlock()

	// create the shutdown channel
	m.shutdown = make(chan bool, 1)

	// lock the database
	err := m.driver.Lock()
	if err != nil {
		return err
	}
	defer m.driver.Unlock()

	// get current version and dirty state
	curVersion, dirty, err := m.driver.GetVersion()
	if err != nil {
		return err
	}

	if dirty {
		return ErrDatabaseDirty
	}

	// get all migrations
	migrations := make(chan *migrationData)
	err = m.GetMigrations(curVersion, version, migrations)
	if err == ErrNoChange {
		m.logger.Printf("no database migration necessary")
		return nil // nothing to do, no error
	}
	if err != nil {
		return err
	}

	// apply all migrations
	return m.applyMigrations(migrations)
}

// GetMigrations fills up a channel with migrations in the background. If the initialization fails, an
// error is returned.
// The migrations channel will be closed by this function.
func (m *migrator) GetMigrations(currentVersion, targetVersion uint64, migrations chan<- *migrationData) error {
	var direction = Up
	if targetVersion < currentVersion {
		direction = Down
	}

	// validate migration versions
	if currentVersion == targetVersion {
		close(migrations)
		return ErrNoChange
	}

	if targetVersion != NoMigrationVersion && !m.isMigrationValid(targetVersion, direction) {
		return fmt.Errorf("invalid target migration version %d", targetVersion)
	}

	// read all migrations
	go func() {
		defer close(migrations)
		var err error

		version := currentVersion // starting target version
		if direction == Up {      // in case we go up, we do not want to apply the current version again
			version, err = m.getNextMigrationVersion(version, direction)
			if err != nil {
				m.logger.Printf("failed to fetch next migration start version: %v", err)
				return
			}
		}

		running := true
		for running {
			if !m.isMigrationValid(version, direction) {
				m.logger.Printf("invalid migration %d, %s", version, direction)
				break
			}

			migration := m.getMigration(version, direction)
			select {
			case <-m.shutdown: // avoid a blocked goroutine by checking if the migrator was shut down
				running = false
				continue
			case migrations <- migration:
			}

			version, err = m.getNextMigrationVersion(version, direction)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				break // no more versions available
			}
			if err != nil {
				m.logger.Printf("failed to fetch next migration version: %v", err)
				break
			}

			// check if all possible migrations are completed
			if direction == Down && version == targetVersion {
				break // reached target version
			}
			if direction == Up && version > targetVersion {
				break // reached target version
			}
		}
	}()

	return nil
}

func (m *migrator) isMigrationValid(version uint64, direction Direction) bool {
	var err error
	var contents io.ReadCloser

	// check if reading the migration contents works
	switch direction {
	case Up:
		contents, _, err = m.source.ReadUp(version)
	case Down:
		contents, _, err = m.source.ReadDown(version)
	}
	if err == nil {
		_ = contents.Close()
	}

	return err == nil
}

func (m *migrator) getNextMigrationVersion(version uint64, direction Direction) (uint64, error) {
	var err error
	var next uint64

	switch direction {
	case Up:
		next, err = m.source.Next(version)
	case Down:
		next, err = m.source.Prev(version)
	}

	return next, err
}

func (m *migrator) getMigration(version uint64, direction Direction) *migrationData {
	var err error
	var contents io.ReadCloser
	var identifier string

	switch direction {
	case Up:
		contents, identifier, err = m.source.ReadUp(version)
	case Down:
		contents, identifier, err = m.source.ReadDown(version)
	}

	targetVersion := version
	if direction == Down {
		targetVersion = targetVersion - 1
	}

	return &migrationData{
		Version:       version,
		TargetVersion: targetVersion,
		Identifier:    identifier,
		Direction:     direction,
		Contents:      contents,
		error:         err,
	}
}

func (m *migrator) applyMigrations(migrations <-chan *migrationData) error {
	defer func() {
		m.shutdown <- true // on error - shutdown migration producer
	}()

	for migration := range migrations {
		// Check if there was an error
		if migration.Error() != nil {
			return migration.Error()
		}

		// Set version with dirty state
		err := m.driver.SetVersion(migration.TargetVersion, true)
		if err != nil {
			return err
		}

		// Apply migration
		err = m.driver.RunMigration(migration.Contents)
		if err != nil {
			_ = migration.Contents.Close()
			return err
		}
		_ = migration.Contents.Close()

		// Remove dirty state
		err = m.driver.SetVersion(migration.TargetVersion, false)
		if err != nil {
			return err
		}

		if m.verbose {
			m.logger.Printf("applied %d, %s (%s)", migration.Version, migration.Direction, migration.Identifier)
		}
	}

	return nil
}

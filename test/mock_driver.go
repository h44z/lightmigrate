package test

import (
	"io"
)

// MockDriver is a mocked driver implementation used for testing.
type MockDriver struct {
	Error   error
	Version uint64
	Dirty   bool
}

// NewMockDriver instantiates a new mocked driver.
func NewMockDriver() (*MockDriver, error) {
	return &MockDriver{}, nil
}

// Close is part of lightmigrate.MigrationDriver interface implementation.
func (m *MockDriver) Close() error {
	return m.Error
}

// Lock is part of lightmigrate.MigrationDriver interface implementation.
func (m *MockDriver) Lock() error {
	return m.Error
}

// Unlock is part of lightmigrate.MigrationDriver interface implementation.
func (m *MockDriver) Unlock() error {
	return m.Error
}

// GetVersion is part of lightmigrate.MigrationDriver interface implementation.
func (m *MockDriver) GetVersion() (version uint64, dirty bool, err error) {
	err = m.Error
	version = m.Version
	dirty = m.Dirty
	return
}

// SetVersion is part of lightmigrate.MigrationDriver interface implementation.
func (m *MockDriver) SetVersion(version uint64, dirty bool) error {
	m.Version = version
	m.Dirty = dirty
	return m.Error
}

// RunMigration is part of lightmigrate.MigrationDriver interface implementation.
func (m *MockDriver) RunMigration(migration io.Reader) error {
	return m.Error
}

// Reset is part of lightmigrate.MigrationDriver interface implementation.
func (m *MockDriver) Reset() error {
	return m.Error
}

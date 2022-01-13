package test

import (
	"io"
)

type MockDriver struct {
	Error   error
	Version uint64
	Dirty   bool
}

func NewMockDriver() (*MockDriver, error) {
	return &MockDriver{}, nil
}

func (m *MockDriver) Close() error {
	return m.Error
}

func (m *MockDriver) Lock() error {
	return m.Error
}

func (m *MockDriver) Unlock() error {
	return m.Error
}

func (m *MockDriver) GetVersion() (version uint64, dirty bool, err error) {
	err = m.Error
	version = m.Version
	dirty = m.Dirty
	return
}

func (m *MockDriver) SetVersion(version uint64, dirty bool) error {
	m.Version = version
	m.Dirty = dirty
	return m.Error
}

func (m *MockDriver) RunMigration(migration io.Reader) error {
	return m.Error
}

func (m *MockDriver) Reset() error {
	return m.Error
}

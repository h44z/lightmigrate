package test

import (
	"bytes"
	"io"
	"io/fs"
	"strconv"
)

type MockSource struct {
	Error      error
	MinVersion uint64
	MaxVersion uint64
	Identifier string
	Contents   []byte
}

func NewAdvancedMockSource(min, max uint64) (*MockSource, error) {
	return &MockSource{
		MinVersion: min,
		MaxVersion: max,
	}, nil
}

func (m *MockSource) Close() error {
	return m.Error
}

func (m *MockSource) First() (version uint64, err error) {
	return m.MinVersion, m.Error
}

func (m *MockSource) Prev(version uint64) (prevVersion uint64, err error) {
	if version > m.MinVersion {
		return version - 1, m.Error
	}

	return 0, &fs.PathError{
		Op:   "prev",
		Path: "/prev",
		Err:  fs.ErrNotExist,
	}
}

func (m *MockSource) Next(version uint64) (nextVersion uint64, err error) {
	if version < m.MaxVersion {
		return version + 1, m.Error
	}

	return 0, &fs.PathError{
		Op:   "next",
		Path: "/next",
		Err:  fs.ErrNotExist,
	}
}

func (m *MockSource) ReadUp(version uint64) (r io.ReadCloser, identifier string, err error) {
	if version < m.MinVersion || version > m.MaxVersion {
		return nil, "", &fs.PathError{
			Op:   "read up for version " + strconv.FormatUint(version, 10),
			Path: "/rup",
			Err:  fs.ErrNotExist,
		}
	}

	return io.NopCloser(bytes.NewReader(m.Contents)), m.Identifier, m.Error
}

func (m *MockSource) ReadDown(version uint64) (r io.ReadCloser, identifier string, err error) {
	if version < m.MinVersion || version > m.MaxVersion {
		return nil, "", &fs.PathError{
			Op:   "read down for version " + strconv.FormatUint(version, 10),
			Path: "/rdown",
			Err:  fs.ErrNotExist,
		}
	}

	return io.NopCloser(bytes.NewReader(m.Contents)), m.Identifier, m.Error
}

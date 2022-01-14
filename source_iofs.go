package lightmigrate

import (
	"errors"
	"io"
	"io/fs"
	"path"
	"strconv"
)

type fsSource struct {
	migrations *migrations

	fsys fs.FS
	path string
}

// NewFsSource returns a new MigrationSource from io/fs#FS and a relative path.
func NewFsSource(fsys fs.FS, basePath string) (MigrationSource, error) {
	f := &fsSource{
		migrations: newMigrations(),
		fsys:       fsys,
		path:       basePath,
	}

	err := f.init()
	if err != nil {
		return nil, err
	}

	return f, nil
}

// init prepares not initialized IoFS instance to read migrations from an
// io/fs#FS instance and a relative path.
func (f *fsSource) init() error {
	entries, err := fs.ReadDir(f.fsys, f.path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		m, err := parseFileName(e.Name())
		if err != nil {
			continue
		}
		file, err := e.Info()
		if err != nil {
			return err
		}
		if !f.migrations.Append(m) {
			return ErrDuplicateMigration{
				migration: *m,
				FileInfo:  file,
			}
		}
	}

	return nil
}

// open a given file path in the filesystem.
func (f *fsSource) open(path string) (fs.File, error) {
	file, err := f.fsys.Open(path)
	if err == nil {
		return file, nil
	}
	// Some non-standard file systems may return errors that don't include the path, that
	// makes debugging harder.
	if !errors.As(err, new(*fs.PathError)) {
		err = &fs.PathError{
			Op:   "open",
			Path: path,
			Err:  err,
		}
	}
	return nil, err
}

// Close is part of source.Driver interface implementation.
// Closes the file system if possible.
func (f *fsSource) Close() error {
	c, ok := f.fsys.(io.Closer)
	if !ok {
		return nil
	}
	return c.Close()
}

// First is part of source.Driver interface implementation.
func (f *fsSource) First() (version uint64, err error) {
	if version, ok := f.migrations.First(); ok {
		return version, nil
	}
	return 0, &fs.PathError{
		Op:   "first",
		Path: f.path,
		Err:  fs.ErrNotExist,
	}
}

// Prev is part of source.Driver interface implementation.
func (f *fsSource) Prev(version uint64) (prevVersion uint64, err error) {
	if version, ok := f.migrations.Prev(version); ok {
		return version, nil
	}
	return 0, &fs.PathError{
		Op:   "prev for version " + strconv.FormatUint(version, 10),
		Path: f.path,
		Err:  fs.ErrNotExist,
	}
}

// Next is part of source.Driver interface implementation.
func (f *fsSource) Next(version uint64) (nextVersion uint64, err error) {
	if version, ok := f.migrations.Next(version); ok {
		return version, nil
	}
	return 0, &fs.PathError{
		Op:   "next for version " + strconv.FormatUint(version, 10),
		Path: f.path,
		Err:  fs.ErrNotExist,
	}
}

// ReadUp is part of source.Driver interface implementation.
func (f *fsSource) ReadUp(version uint64) (r io.ReadCloser, identifier string, err error) {
	if m, ok := f.migrations.Up(version); ok {
		body, err := f.open(path.Join(f.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return body, m.Identifier, nil
	}
	return nil, "", &fs.PathError{
		Op:   "read up for version " + strconv.FormatUint(version, 10),
		Path: f.path,
		Err:  fs.ErrNotExist,
	}
}

// ReadDown is part of source.Driver interface implementation.
func (f *fsSource) ReadDown(version uint64) (r io.ReadCloser, identifier string, err error) {
	if m, ok := f.migrations.Down(version); ok {
		body, err := f.open(path.Join(f.path, m.Raw))
		if err != nil {
			return nil, "", err
		}
		return body, m.Identifier, nil
	}
	return nil, "", &fs.PathError{
		Op:   "read down for version " + strconv.FormatUint(version, 10),
		Path: f.path,
		Err:  fs.ErrNotExist,
	}
}

package lightmigrate

import (
	"bytes"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type closeableFs struct{}

func (c closeableFs) Open(name string) (fs.File, error) { return nil, nil }
func (c closeableFs) Close() error                      { return nil }

func getTestSource(t *testing.T, folder string) *fsSource {
	fsys := os.DirFS("test")
	source, err := NewFsSource(fsys, folder)
	if err != nil {
		t.Fatalf("unable to setup source: %v", err)
	}
	return source.(*fsSource)
}

func TestNewFsSource(t *testing.T) {
	fsys := os.DirFS("test")
	source, err := NewFsSource(fsys, "sample-migrations")
	if err != nil {
		t.Fatalf("unable to setup source: %v", err)
	}
	defer source.Close()
}

func TestNewFsSource_WrongDirectory(t *testing.T) {
	fsys := os.DirFS("test_nodir")
	_, err := NewFsSource(fsys, "sample-migrations")
	if err == nil {
		t.Fatal("expected an error")
	}
}

func Test_fsSource_Close_NothingToClose(t *testing.T) {
	fsys := os.DirFS("test")
	f := &fsSource{
		fsys: fsys,
	}

	err := f.Close()
	if err != nil {
		t.Fatalf("unable to close source: %v", err)
	}
}

func Test_fsSource_Close(t *testing.T) {
	f := &fsSource{
		fsys: closeableFs{},
	}

	err := f.Close()
	if err != nil {
		t.Fatalf("unable to close source: %v", err)
	}
}

func Test_fsSource_First(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	version, err := s.First()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != uint64(1) {
		t.Fatalf("expected first version to be 1, got: %d", version)
	}
}

func Test_fsSource_First_NoMigrations(t *testing.T) {
	s := getTestSource(t, "no-migrations")
	version, err := s.First()
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got: %v", err)
	}
	if version != 0 {
		t.Fatalf("expected first version to be 0, got: %d", version)
	}

}

func Test_fsSource_Next_FromStart(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	version, err := s.Next(NoMigrationVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != uint64(1) {
		t.Fatalf("expected first version to be 1, got: %d", version)
	}
}

func Test_fsSource_Next(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	version, err := s.Next(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != uint64(2) {
		t.Fatalf("expected first version to be 2, got: %d", version)
	}
}

func Test_fsSource_Next_NoNext(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	version, err := s.Next(3)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got: %v", err)
	}
	if version != 0 {
		t.Fatalf("expected first version to be 0, got: %d", version)
	}
}

func Test_fsSource_Prev_FromStart(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	version, err := s.Prev(NoMigrationVersion)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got: %v", err)
	}
	if version != 0 {
		t.Fatalf("expected first version to be 0, got: %d", version)
	}
}

func Test_fsSource_Prev(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	version, err := s.Prev(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != uint64(1) {
		t.Fatalf("expected first version to be 2, got: %d", version)
	}
}

func Test_fsSource_Prev_InvalidFrom(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	version, err := s.Prev(5)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got: %v", err)
	}
	if version != 0 {
		t.Fatalf("expected first version to be 0, got: %d", version)
	}
}

func Test_fsSource_ReadDown(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	down, identifier, err := s.ReadDown(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identifier != "some-text" {
		t.Fatalf("expected identifier to be some-text, got: %s", identifier)
	}
	defer down.Close()
	contents, _ := ioutil.ReadAll(down)
	if bytes.Compare(contents, []byte("{\"1\": \"down\"}")) != 0 {
		t.Fatalf("unexpected contents, got: %s", contents)
	}
}

func Test_fsSource_ReadDown_NoMigration(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	_, _, err := s.ReadDown(NoMigrationVersion)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got: %v", err)
	}
}

func Test_fsSource_ReadDown_FileError(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	s.migrations.Append(&migration{
		Version:    4,
		Identifier: "invalid",
		Direction:  "down",
		Raw:        "no_such_file",
	})
	_, _, err := s.ReadDown(4)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got: %v", err)
	}
}

func Test_fsSource_ReadUp(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	up, identifier, err := s.ReadUp(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identifier != "some-text" {
		t.Fatalf("expected identifier to be some-text, got: %s", identifier)
	}
	defer up.Close()
	contents, _ := ioutil.ReadAll(up)
	if bytes.Compare(contents, []byte("{\"1\": \"up\"}")) != 0 {
		t.Fatalf("unexpected contents, got: %s", contents)
	}
}

func Test_fsSource_ReadUp_NoMigration(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	_, _, err := s.ReadUp(4)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got: %v", err)
	}
}

func Test_fsSource_ReadUp_FileError(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	s.migrations.Append(&migration{
		Version:    4,
		Identifier: "invalid",
		Direction:  "up",
		Raw:        "no_such_file",
	})
	_, _, err := s.ReadUp(4)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got: %v", err)
	}
}

func Test_fsSource_open(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	f, err := s.open(path.Join(s.path, "002_another_text.up.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer f.Close()
	contents, _ := ioutil.ReadAll(f)
	if bytes.Compare(contents, []byte("{\"2\": \"up\"}")) != 0 {
		t.Fatalf("unexpected contents, got: %s", contents)
	}
}

func Test_fsSource_open_NoSuchFile(t *testing.T) {
	s := getTestSource(t, "sample-migrations")
	_, err := s.open("invalid path")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got: %v", err)
	}
}

func Test_fsSource_init(t *testing.T) {
	fsys := os.DirFS("test")
	f := &fsSource{
		migrations: newMigrations(),
		fsys:       fsys,
		path:       "sample-migrations",
	}

	err := f.init()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.migrations.index) != 3 {
		t.Fatalf("expected migrations index to have length of 3, got: %d", len(f.migrations.index))
	}
}

func Test_fsSource_init_NoMigrations(t *testing.T) {
	fsys := os.DirFS("test")
	f := &fsSource{
		migrations: newMigrations(),
		fsys:       fsys,
		path:       "no-migrations",
	}

	err := f.init()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.migrations.index) != 0 {
		t.Fatalf("expected migrations index to have length of 0, got: %d", len(f.migrations.index))
	}
}

func Test_fsSource_init_DuplicateMigrations(t *testing.T) {
	fsys := os.DirFS("test")
	f := &fsSource{
		migrations: newMigrations(),
		fsys:       fsys,
		path:       "duplicate-migrations",
	}

	err := f.init()
	if err == nil {
		t.Fatalf("expected ErrDuplicateMigration, got: %v", err)
	}
}

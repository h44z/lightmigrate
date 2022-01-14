package lightmigrate

import (
	"os"
	"testing"

	"github.com/h44z/lightmigrate/test"
)

var sampleFileRoot = "test"
var sampleFilePath = "sample-migrations"

func TestMockedMigratorUp(t *testing.T) {
	fsys := os.DirFS(sampleFileRoot)

	source, err := NewFsSource(fsys, sampleFilePath)
	if err != nil {
		t.Fatalf("unable to setup source: %v", err)
	}
	defer source.Close()

	driver, err := test.NewMockDriver()
	if err != nil {
		t.Fatalf("unable to setup driver: %v", err)
	}
	defer driver.Close()

	migrator, err := NewMigrator(source, driver, WithVerboseLogging(true))
	if err != nil {
		t.Fatalf("unable to setup migrator: %v", err)
	}

	// Invalid version
	err = migrator.Migrate(4)
	if err == nil {
		t.Fatal("expected a migration error")
	}
	if v, _, _ := driver.GetVersion(); v != 0 {
		t.Fatal("expected driver to be on migration version 0")
	}

	// Valid version
	err = migrator.Migrate(2)
	if err != nil {
		t.Fatalf("expected no migration error: %v", err)
	}
	if v, _, _ := driver.GetVersion(); v != 2 {
		t.Fatal("expected driver to be on migration version 2")
	}

	// Up one more
	err = migrator.Migrate(3)
	if err != nil {
		t.Fatalf("expected no migration error: %v", err)
	}
	if v, _, _ := driver.GetVersion(); v != 3 {
		t.Fatal("expected driver to be on migration version 3")
	}

	// Same again
	err = migrator.Migrate(3)
	if err != nil {
		t.Fatalf("expected no migration error: %v", err)
	}
	if v, _, _ := driver.GetVersion(); v != 3 {
		t.Fatal("expected driver to be on migration version 3")
	}
}

func TestMockedMigratorDown(t *testing.T) {
	fsys := os.DirFS(sampleFileRoot)

	source, err := NewFsSource(fsys, sampleFilePath)
	if err != nil {
		t.Fatalf("unable to setup source: %v", err)
	}
	defer source.Close()

	driver, err := test.NewMockDriver()
	if err != nil {
		t.Fatalf("unable to setup driver: %v", err)
	}
	driver.Version = 2
	defer driver.Close()

	migrator, err := NewMigrator(source, driver, WithVerboseLogging(true))
	if err != nil {
		t.Fatalf("unable to setup migrator: %v", err)
	}

	// Rollback to 2 (no change)
	err = migrator.Migrate(2)
	if err != nil {
		t.Fatalf("expected no migration error: %v", err)
	}
	if v, _, _ := driver.GetVersion(); v != 2 {
		t.Fatal("expected driver to be on migration version 2")
	}

	// Rollback to 1
	err = migrator.Migrate(1)
	if err != nil {
		t.Fatalf("expected no migration error: %v", err)
	}
	if v, _, _ := driver.GetVersion(); v != 1 {
		t.Fatal("expected driver to be on migration version 1")
	}

	// Rollback to 0
	err = migrator.Migrate(0) // 0 = remove all migrations
	if err != nil {
		t.Fatalf("expected no migration error: %v", err)
	}
	if v, _, _ := driver.GetVersion(); v != 0 {
		t.Fatal("expected driver to be on migration version 0")
	}

	// Rollback again
	err = migrator.Migrate(0) // 0 = remove all migrations
	if err != nil {
		t.Fatalf("expected no migration error: %v", err)
	}
	if v, _, _ := driver.GetVersion(); v != 0 {
		t.Fatal("expected driver to be on migration version 0")
	}
}

func TestMockedMigratorDirty(t *testing.T) {
	fsys := os.DirFS(sampleFileRoot)

	source, err := NewFsSource(fsys, sampleFilePath)
	if err != nil {
		t.Fatalf("unable to setup source: %v", err)
	}
	defer source.Close()

	driver, err := test.NewMockDriver()
	if err != nil {
		t.Fatalf("unable to setup driver: %v", err)
	}
	driver.Version = 2
	driver.Dirty = true
	defer driver.Close()

	migrator, err := NewMigrator(source, driver, WithVerboseLogging(true))
	if err != nil {
		t.Fatalf("unable to setup migrator: %v", err)
	}

	// Rollback to 1
	err = migrator.Migrate(1)
	if err == nil {
		t.Fatal("expected migration error")
	}
	t.Log(err)
	if v, _, _ := driver.GetVersion(); v != 2 {
		t.Fatal("expected driver to be on migration version 2")
	}
}

package lightmigrate

import (
	"bytes"
	"errors"
	"github.com/h44z/lightmigrate/test"
	"io"
	"log"
	"reflect"
	"testing"
)

func getTestMigrator() *migrator {
	d, _ := test.NewMockDriver()
	s, _ := test.NewMockSource(1, 2)
	m, _ := NewMigrator(s, d)

	return m.(*migrator)
}

func TestNewMigrator(t *testing.T) {
	_, err := NewMigrator(nil, nil, WithLogger(log.Default()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWithLogger(t *testing.T) {
	m := &migrator{}

	WithLogger(log.Default())(m)
	if m.logger != log.Default() {
		t.Fatalf("failed to set logger")
	}
}

func TestWithVerboseLogging(t *testing.T) {
	m := &migrator{}

	WithVerboseLogging(true)(m)
	if m.verbose != true {
		t.Fatalf("failed to set verbose flag")
	}
}

func Test_migrationData_Error(t *testing.T) {
	m := migrationData{error: ErrNoChange}

	if got := m.Error(); got != ErrNoChange {
		t.Fatalf("expected error: %v, got: %v", ErrNoChange, got)
	}
}

func Test_migrator_GetMigrations(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 2)
	m.source = s

	migrations := make(chan *migrationData)
	err := m.GetMigrations(NoMigrationVersion, 1, migrations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions := make([]uint64, 0)
	wantVersions := []uint64{1}
	for mig := range migrations {
		if mig.Error() != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		versions = append(versions, mig.Version)
	}
	if !reflect.DeepEqual(versions, wantVersions) {
		t.Fatalf("expected versions: %v, got: %v", wantVersions, versions)
	}
}

func Test_migrator_GetMigrations_NoChange(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 2)
	m.source = s

	migrations := make(chan *migrationData)
	err := m.GetMigrations(1, 1, migrations)
	if err != ErrNoChange {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_migrator_GetMigrations_Shutdown(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	m.source = s
	m.shutdown = make(chan bool, 1)

	migrations := make(chan *migrationData)
	err := m.GetMigrations(NoMigrationVersion, 3, migrations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions := make([]uint64, 0)
	wantVersions := []uint64{1}
	for mig := range migrations {
		if mig.Error() != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		versions = append(versions, mig.Version)

		m.shutdown <- true // shutdown after first migration
		break
	}

	for mig := range migrations {
		if mig.Error() != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		versions = append(versions, mig.Version)
	}
	if !reflect.DeepEqual(versions, wantVersions) {
		t.Fatalf("expected versions: %v, got: %v", wantVersions, versions)
	}
}

func Test_migrator_GetMigrations_MultiUp(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	m.source = s

	migrations := make(chan *migrationData)
	err := m.GetMigrations(1, 3, migrations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions := make([]uint64, 0)
	wantVersions := []uint64{2, 3}
	for mig := range migrations {
		if mig.Error() != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		versions = append(versions, mig.Version)
	}
	if !reflect.DeepEqual(versions, wantVersions) {
		t.Fatalf("expected versions: %v, got: %v", wantVersions, versions)
	}
}

func Test_migrator_GetMigrations_MultiUpZero(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	m.source = s

	migrations := make(chan *migrationData)
	err := m.GetMigrations(NoMigrationVersion, 3, migrations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions := make([]uint64, 0)
	wantVersions := []uint64{1, 2, 3}
	for mig := range migrations {
		if mig.Error() != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		versions = append(versions, mig.Version)
	}
	if !reflect.DeepEqual(versions, wantVersions) {
		t.Fatalf("expected versions: %v, got: %v", wantVersions, versions)
	}
}

func Test_migrator_GetMigrations_MultiDown(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	m.source = s

	migrations := make(chan *migrationData)
	err := m.GetMigrations(3, 1, migrations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions := make([]uint64, 0)
	wantVersions := []uint64{3, 2}
	for mig := range migrations {
		if mig.Error() != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		versions = append(versions, mig.Version)
	}
	if !reflect.DeepEqual(versions, wantVersions) {
		t.Fatalf("expected versions: %v, got: %v", wantVersions, versions)
	}
}

func Test_migrator_GetMigrations_MultiDownZero(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	m.source = s

	migrations := make(chan *migrationData)
	err := m.GetMigrations(3, NoMigrationVersion, migrations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions := make([]uint64, 0)
	wantVersions := []uint64{3, 2, 1}
	for mig := range migrations {
		if mig.Error() != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		versions = append(versions, mig.Version)
	}
	if !reflect.DeepEqual(versions, wantVersions) {
		t.Fatalf("expected versions: %v, got: %v", wantVersions, versions)
	}
}

func Test_migrator_GetMigrations_WrongTarget(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 2)
	m.source = s

	migrations := make(chan *migrationData)
	err := m.GetMigrations(NoMigrationVersion, 3, migrations)
	if err == nil {
		t.Fatalf("expected error: %v", err)
	}
}

func Test_migrator_Migrate(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	m.source = s

	err := m.Migrate(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_migrator_Migrate_NoChange(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	m.source = s

	err := m.Migrate(NoMigrationVersion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_migrator_Migrate_DriverError(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	m.source = s
	m.driver.(*test.MockDriver).Error = errors.New("lockerror")

	err := m.Migrate(2)
	if err == nil {
		t.Fatalf("expected error, got nil: %v", err)
	}
}

func Test_migrator_Migrate_DriverDirty(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	m.source = s
	m.driver.(*test.MockDriver).Dirty = true

	err := m.Migrate(2)
	if err != ErrDatabaseDirty {
		t.Fatalf("expected ErrDatabaseDirty error, got: %v", err)
	}
}

func Test_migrator_Migrate_SourceError(t *testing.T) {
	m := getTestMigrator()
	s, _ := test.NewMockSource(1, 3)
	s.Error = errors.New("sourcerror")
	m.source = s

	err := m.Migrate(2)
	if err == nil {
		t.Fatalf("expected error, got nil: %v", err)
	}
}

func Test_migrator_applyMigrations(t *testing.T) {
	m := getTestMigrator()
	m.shutdown = make(chan bool, 1)

	migrations := make(chan *migrationData, 2)
	migrations <- &migrationData{Version: 1, Contents: io.NopCloser(bytes.NewReader([]byte("test1")))}
	migrations <- &migrationData{Version: 2, Contents: io.NopCloser(bytes.NewReader([]byte("test2")))}
	close(migrations)

	err := m.applyMigrations(migrations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_migrator_applyMigrations_Verbose(t *testing.T) {
	m := getTestMigrator()
	m.verbose = true
	m.shutdown = make(chan bool, 1)

	migrations := make(chan *migrationData, 2)
	migrations <- &migrationData{Version: 1, Contents: io.NopCloser(bytes.NewReader([]byte("test1")))}
	migrations <- &migrationData{Version: 2, Contents: io.NopCloser(bytes.NewReader([]byte("test2")))}
	close(migrations)

	err := m.applyMigrations(migrations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_migrator_applyMigrations_MigrationError(t *testing.T) {
	m := getTestMigrator()
	m.shutdown = make(chan bool, 1)

	migrations := make(chan *migrationData, 1)
	migrations <- &migrationData{error: ErrVersionNotAllowed}

	err := m.applyMigrations(migrations)
	if err != ErrVersionNotAllowed {
		t.Fatalf("unexpected error: %v", err)
	}
	if s := <-m.shutdown; !s {
		t.Fatalf("not shut down: %t", s)
	}
}

func Test_migrator_getMigration_Up(t *testing.T) {
	m := getTestMigrator()

	got := m.getMigration(1, Up)
	if got.Version != 1 {
		t.Fatalf("unexpected version: %d", got.Version)
	}
	if got.TargetVersion != 1 {
		t.Fatalf("unexpected target version: %d", got.Version)
	}
}

func Test_migrator_getMigration_Down(t *testing.T) {
	m := getTestMigrator()

	got := m.getMigration(1, Down)
	if got.Version != 1 {
		t.Fatalf("unexpected version: %d", got.Version)
	}
	if got.TargetVersion != NoMigrationVersion {
		t.Fatalf("unexpected target version: %d", got.Version)
	}
}

func Test_migrator_getNextMigrationVersion_Up(t *testing.T) {
	m := getTestMigrator()

	got, err := m.getNextMigrationVersion(1, Up)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 2 {
		t.Fatalf("unexpected version: %d", got)
	}
}

func Test_migrator_getNextMigrationVersion_Down(t *testing.T) {
	m := getTestMigrator()

	got, err := m.getNextMigrationVersion(2, Down)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 1 {
		t.Fatalf("unexpected version: %d", got)
	}
}

func Test_migrator_isMigrationValid(t *testing.T) {
	m := getTestMigrator()

	if got := m.isMigrationValid(1, Up); !got {
		t.Fatalf("expected migration to be valid, got %t", got)
	}

	if got := m.isMigrationValid(1, Down); !got {
		t.Fatalf("expected migration to be valid, got %t", got)
	}
}

func Test_migrator_isMigrationValid_Invalid(t *testing.T) {
	m := getTestMigrator()
	m.source.(*test.MockSource).Error = errors.New("errmsg")

	if got := m.isMigrationValid(1, Up); got {
		t.Fatalf("expected migration to be invalid, got %t", got)
	}

	if got := m.isMigrationValid(1, Down); got {
		t.Fatalf("expected migration to be invalid, got %t", got)
	}
}

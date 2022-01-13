package lightmigrate

import (
	"testing"
)

func TestNewMigrations(t *testing.T) {
	m := newMigrations()
	if m.migrations == nil {
		t.Fatalf("migrations uninitialized")
	}
	if m.index == nil {
		t.Fatalf("index uninitialized")
	}
}

func Test_migrations_Append(t *testing.T) {
	m := newMigrations()

	ok := m.Append(&migration{})
	if !ok {
		t.Fatalf("failed to append")
	}
	if len(m.index) != 1 {
		t.Fatalf("expected index to have length of 1, got: %d", len(m.index))
	}
}

func Test_migrations_Append_Multi(t *testing.T) {
	m := newMigrations()

	ok := m.Append(&migration{Version: 1})
	if !ok {
		t.Fatalf("failed to append")
	}
	ok = m.Append(&migration{Version: 2})
	if !ok {
		t.Fatalf("failed to append")
	}
	if len(m.index) != 2 {
		t.Fatalf("expected index to have length of 2, got: %d", len(m.index))
	}
}

func Test_migrations_Append_Duplicate(t *testing.T) {
	m := newMigrations()

	ok := m.Append(&migration{Version: 1})
	if !ok {
		t.Fatalf("failed to append")
	}
	ok = m.Append(&migration{Version: 1})
	if ok {
		t.Fatalf("append should not have worked")
	}
	if len(m.index) != 1 {
		t.Fatalf("expected index to have length of 1, got: %d", len(m.index))
	}
}

func Test_migrations_Append_Nil(t *testing.T) {
	m := newMigrations()

	ok := m.Append(nil)
	if ok {
		t.Fatalf("append should not have worked")
	}
}

func Test_migrations_Down(t *testing.T) {
	m := newMigrations()
	mig := &migration{Version: 1, Direction: Down}
	m.Append(mig)

	down, ok := m.Down(1)
	if !ok {
		t.Fatalf("failed to get migration")
	}

	if down != mig {

		t.Fatalf("invalid migration, got: %v", down)
	}
}

func Test_migrations_Down_NoMigration(t *testing.T) {
	m := newMigrations()
	mig := &migration{Version: 1, Direction: Up}
	m.Append(mig)

	_, ok := m.Down(1)
	if ok {
		t.Fatalf("unexpected migration")
	}
}

func Test_migrations_First(t *testing.T) {
	m := newMigrations()
	m.Append(&migration{Version: 2, Direction: Down})
	m.Append(&migration{Version: 1, Direction: Down})
	m.Append(&migration{Version: 3, Direction: Down})

	version, ok := m.First()
	if !ok {
		t.Fatalf("unexpected error")
	}
	if version != 1 {
		t.Fatalf("unexpected migration version: %d", version)
	}
}

func Test_migrations_First_NoMigration(t *testing.T) {
	m := newMigrations()

	version, ok := m.First()
	if ok {
		t.Fatalf("unexpected migration")
	}
	if version != NoMigrationVersion {
		t.Fatalf("unexpected migration version: %d", version)
	}
}

func Test_migrations_Next(t *testing.T) {
	m := newMigrations()
	m.Append(&migration{Version: 2, Direction: Down})
	m.Append(&migration{Version: 1, Direction: Down})
	m.Append(&migration{Version: 3, Direction: Down})

	version, ok := m.Next(NoMigrationVersion)
	if !ok {
		t.Fatalf("unexpected error")
	}
	if version != 1 {
		t.Fatalf("unexpected migration version: %d", version)
	}
}

func Test_migrations_Next_FromStart(t *testing.T) {
	m := newMigrations()
	m.Append(&migration{Version: 2, Direction: Down})
	m.Append(&migration{Version: 1, Direction: Down})
	m.Append(&migration{Version: 3, Direction: Down})

	version, ok := m.Next(2)
	if !ok {
		t.Fatalf("unexpected error")
	}
	if version != 3 {
		t.Fatalf("unexpected migration version: %d", version)
	}
}

func Test_migrations_Next_NoNext(t *testing.T) {
	m := newMigrations()
	m.Append(&migration{Version: 2, Direction: Down})
	m.Append(&migration{Version: 1, Direction: Down})
	m.Append(&migration{Version: 3, Direction: Down})

	version, ok := m.Next(3)
	if ok {
		t.Fatalf("unexpected migration")
	}
	if version != 0 {
		t.Fatalf("unexpected migration version: %d", version)
	}
}

func Test_migrations_Prev(t *testing.T) {
	m := newMigrations()
	m.Append(&migration{Version: 2, Direction: Down})
	m.Append(&migration{Version: 1, Direction: Down})
	m.Append(&migration{Version: 3, Direction: Down})

	version, ok := m.Prev(2)
	if !ok {
		t.Fatalf("unexpected error")
	}
	if version != 1 {
		t.Fatalf("unexpected migration version: %d", version)
	}
}

func Test_migrations_Prev_NoPrev(t *testing.T) {
	m := newMigrations()
	m.Append(&migration{Version: 2, Direction: Down})
	m.Append(&migration{Version: 1, Direction: Down})
	m.Append(&migration{Version: 3, Direction: Down})

	version, ok := m.Prev(NoMigrationVersion)
	if ok {
		t.Fatalf("unexpected migration")
	}
	if version != 0 {
		t.Fatalf("unexpected migration version: %d", version)
	}
}

func Test_migrations_Up(t *testing.T) {
	m := newMigrations()
	mig := &migration{Version: 1, Direction: Up}
	m.Append(mig)

	up, ok := m.Up(1)
	if !ok {
		t.Fatalf("failed to get migration")
	}

	if up != mig {

		t.Fatalf("invalid migration, got: %v", up)
	}
}

func Test_migrations_Up_NoMigration(t *testing.T) {
	m := newMigrations()
	mig := &migration{Version: 1, Direction: Down}
	m.Append(mig)

	_, ok := m.Up(1)
	if ok {
		t.Fatalf("unexpected migration")
	}
}

func Test_migrations_buildIndex(t *testing.T) {
	m := newMigrations()

	ok := m.Append(&migration{Version: 1})
	if !ok {
		t.Fatalf("failed to append")
	}
	m.buildIndex()
	if len(m.index) != 1 {
		t.Fatalf("expected index to have length of 1, got: %d", len(m.index))
	}
}

func Test_migrations_findPos(t *testing.T) {
	m := newMigrations()
	m.Append(&migration{Version: 1})

	pos := m.findPos(1)
	if pos != 0 {
		t.Fatalf("pos should be 0, got: %d", pos)
	}
}

func Test_migrations_findPos_Invalid(t *testing.T) {
	m := newMigrations()
	m.Append(&migration{Version: 1})

	pos := m.findPos(2)
	if pos != -1 {
		t.Fatalf("pos should be -1, got: %d", pos)
	}
}

func Test_migrations_findPos_Empty(t *testing.T) {
	m := newMigrations()
	pos := m.findPos(1)
	if pos != -1 {
		t.Fatalf("pos should be -1, got: %d", pos)
	}
}

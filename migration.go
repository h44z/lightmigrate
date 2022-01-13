package lightmigrate

import "sort"

// Direction is either up or down.
type Direction string

const (
	// Down direction is used when migrations should be reverted.
	Down Direction = "down"
	// Up direction is the default direction for migrations.
	Up Direction = "up"
)

// migrations wraps migration and has an internal index
// to keep track of migration order.
type migrations struct {
	// Store the order of the versions
	index []uint64
	// Store migrations for each version
	migrations map[uint64]map[Direction]*migration
}

// migration meta object
type migration struct {
	// Version is the version of this migration.
	Version uint64

	// Identifier can be any string that helps identifying
	// this migration in the source.
	Identifier string

	// Direction is either Up or Down.
	Direction Direction

	// Raw holds the raw location path to this migration in source.
	// ReadUp and ReadDown will use this.
	Raw string
}

func newMigrations() *migrations {
	return &migrations{
		index:      make([]uint64, 0),
		migrations: make(map[uint64]map[Direction]*migration),
	}
}

func (i *migrations) Append(m *migration) (ok bool) {
	if m == nil {
		return false
	}

	if i.migrations[m.Version] == nil {
		i.migrations[m.Version] = make(map[Direction]*migration)
	}

	// reject duplicate versions
	if _, dup := i.migrations[m.Version][m.Direction]; dup {
		return false
	}

	i.migrations[m.Version][m.Direction] = m
	i.buildIndex()

	return true
}

func (i *migrations) buildIndex() {
	i.index = make([]uint64, 0, len(i.migrations))
	for version := range i.migrations {
		i.index = append(i.index, version)
	}
	sort.Slice(i.index, func(x, y int) bool {
		return i.index[x] < i.index[y]
	})
}

func (i *migrations) First() (version uint64, ok bool) {
	if len(i.index) == 0 {
		return 0, false
	}
	return i.index[0], true
}

func (i *migrations) Prev(version uint64) (prevVersion uint64, ok bool) {
	pos := i.findPos(version)
	if pos >= 1 && len(i.index) > pos-1 {
		return i.index[pos-1], true
	}
	return 0, false
}

func (i *migrations) Next(version uint64) (nextVersion uint64, ok bool) {
	pos := i.findPos(version)
	if len(i.index) > pos+1 {
		return i.index[pos+1], true
	}
	return 0, false
}

func (i *migrations) Up(version uint64) (m *migration, ok bool) {
	if _, ok := i.migrations[version]; ok {
		if mx, ok := i.migrations[version][Up]; ok {
			return mx, true
		}
	}
	return nil, false
}

func (i *migrations) Down(version uint64) (m *migration, ok bool) {
	if _, ok := i.migrations[version]; ok {
		if mx, ok := i.migrations[version][Down]; ok {
			return mx, true
		}
	}
	return nil, false
}

func (i *migrations) findPos(version uint64) int {
	if len(i.index) > 0 {
		ix := sort.Search(len(i.index), func(j int) bool { return i.index[j] >= version })
		if ix < len(i.index) && i.index[ix] == version {
			return ix
		}
	}
	return -1
}

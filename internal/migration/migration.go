package migration

// TODO: Refactor to make it a generically usable module.
// TODO: Move to its own module (mikrowezel/migration)

type (
	Migrator struct {
		up   []Migration
		down []Rollback
	}

	proc struct {
		order    int
		funct    func() error
		executed bool
		err      error
	}

	Migration struct {
		proc
	}

	Rollback struct {
		proc
	}
)

func init() {
}

func (m *Migrator) CreateDb() error {
	return nil
}

func (m *Migrator) DropDb() error {
	return nil
}

func (m *Migrator) AddUp(mg Migration) {

}

func (m *Migrator) AddDown(rb Rollback) {

}

func (m *Migrator) MigrateAll() {

}

func (m *Migrator) RollbackAll() {

}

func (m *Migrator) MigrateThis(mg Migration) error {
	return nil
}

func (m *Migrator) RollbackThis(r Rollback) error {
	return nil
}

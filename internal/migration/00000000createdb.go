package migration

import (
	"fmt"
)

func init() {
}

// Up00000000 migration
func (m *migrator) Up00000000() error {
	tx := m.getTx()

	st := fmt.Sprintf(`
		CREATE DATABASE %s;
	`, testDb)

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Down00000000 rollback
func (m *migrator) Down00000000() error {
	tx := m.getTx()

	st := fmt.Sprintf(`
		DROP DATABASE %s;
	`, testDb)

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return tx.Commit()
}

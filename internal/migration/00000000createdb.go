package migration

import (
	"fmt"
)

func init() {
}

// Up00000000 migration
func (t transaction) Up00000000() error {
	tx := t.getTx()

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
func (t transaction) Down00000000() error {
	tx := t.getTx()

	st := fmt.Sprintf(`
		DROP DATABASE %s;
	`, testDb)

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return tx.Commit()
}

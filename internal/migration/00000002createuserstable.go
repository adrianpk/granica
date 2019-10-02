package migration

import (
	"database/sql"
	"fmt"
)

// Up00000002 migration
func Up00000002(tx *sql.Tx) error {
	st := fmt.Sprintf(`
		CREATE DATABASE %s;
	`, testDb)

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}
	return nil
}

// Down00000002 migration
func Down00000002(tx *sql.Tx) error {
	st := fmt.Sprintf(`
		DROP DATABASE %s;
	`, testDb)

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}
	return nil
}

package migration

import (
	"database/sql"
	"fmt"
)

func init() {
}

// Up00000001 migration
func Up00000001(tx *sql.Tx) error {
	st := fmt.Sprintf(`
		CREATE DATABASE %s;
	`, testDb)

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}
	return nil
}

// 00000002 migration
func Down00000001(tx *sql.Tx) error {
	st := fmt.Sprintf(`
		DROP DATABASE %s;
	`, testDb)

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}
	return nil
}

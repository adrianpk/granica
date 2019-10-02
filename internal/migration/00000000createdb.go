package migration

import (
	"fmt"
)

const (
	name0 = "create_database"
)

// Up00000000 migration
func (m *migrator) Up00000000() procResult {
	tx := m.getTx()

	st := fmt.Sprintf(`
		CREATE DATABASE %s;
	`, testDb)

	_, err := tx.Exec(st)
	if err != nil {
		return m.makeProcResult(tx, name0, err)
	}

	st = `CREATE TABLE migrations
		id UUID PRIMARY KEY,
		name VARCHAR(32),
 		is_applied BOOLEAN,
		created_at TIMESTAMP;
	);`

	_, err = tx.Exec(st)
	if err != nil {
		return m.makeProcResult(tx, name0, err)
	}

	return m.makeProcResult(tx, name0, tx.Commit())
}

// Down00000000 rollback
func (m *migrator) Down00000000() procResult {
	tx := m.getTx()

	st := fmt.Sprintf(`
		DROP DATABASE %s;
	`, testDb)

	_, err := tx.Exec(st)
	if err != nil {
		return m.makeProcResult(tx, name0, err)
	}

	return m.makeProcResult(tx, name0, tx.Commit())
}

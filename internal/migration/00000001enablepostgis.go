package migration

const (
	name1 = "enable_postgis"
)

// Up00000001 migration
func (m *migrator) Up00000001() procResult {
	tx := m.getTx()

	st := `CREATE EXTENSION IF NOT EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return m.makeProcResult(tx, name1, err)
	}

	return m.makeProcResult(tx, name1, tx.Commit())
}

// Down00000001 rollback
func (m *migrator) Down00000001() procResult {
	tx := m.getTx()

	st := `DROP EXTENSION IF EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return m.makeProcResult(tx, name1, err)
	}

	return m.makeProcResult(tx, name1, tx.Commit())
}

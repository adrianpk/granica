package migration

func init() {
}

// Up00000001 migration
func (m *migrator) Up00000001() error {
	tx := m.getTx()

	st := `CREATE EXTENSION IF NOT EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Down00000001 rollback
func (m *migrator) Down00000001() error {
	tx := m.getTx()

	st := `DROP EXTENSION IF NOT EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return tx.Commit()
}

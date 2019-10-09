package migration

// EnablePostgis migration
func (m *mig) EnablePostgis() error {
	tx := m.GetTx()

	st := `CREATE EXTENSION IF NOT EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// DropPostgis rollback
func (m *mig) DropPostgis() error {
	tx := m.GetTx()

	st := `DROP EXTENSION IF EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

package migration

const (
	name1 = "enable_postgis"
)

// Up00000001 migration
func (m *mig) Up00000001() (string, error) {
	tx := m.GetTx()

	st := `CREATE EXTENSION IF NOT EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return name1, err
	}

	return name1, nil
}

// Down00000001 rollback
func (m *mig) Down00000001() (string, error) {
	tx := m.GetTx()

	st := `DROP EXTENSION IF EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return name1, err
	}

	return name1, nil
}

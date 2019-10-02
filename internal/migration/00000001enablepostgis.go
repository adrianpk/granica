package migration

func init() {
}

// Up00000001 migration
func (t *transaction) Up00000001() error {
	tx := t.getTx()

	st := `CREATE EXTENSION IF NOT EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Down00000001 rollback
func (t *transaction) Down00000001() error {
	tx := t.getTx()

	st := `DROP EXTENSION IF NOT EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return tx.Commit()
}

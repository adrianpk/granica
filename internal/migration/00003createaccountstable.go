package migration

import "log"

// CreateAccountsTable migration
func (m *mig) CreateAccountsTable() error {
	tx := m.GetTx()

	st := `CREATE TABLE accounts
	(
		id UUID PRIMARY KEY,
		tenant_id VARCHAR(128),
		name VARCHAR(64),
		slug VARCHAR(36) UNIQUE,
		owner_id UUID,
		parent_id UUID,
	  account_type VARCHAR(36) UNIQUE,
		email VARCHAR(255) UNIQUE,
		shown_name VARCHAR(128)
	);`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	st = `
		ALTER TABLE accounts
		ADD COLUMN geolocation geography (Point,4326),
		ADD COLUMN locale VARCHAR(32),
		ADD COLUMN base_tz VARCHAR(2),
		ADD COLUMN current_tz VARCHAR(2),
		ADD COLUMN starts_at TIMESTAMP,
		ADD COLUMN ends_at TIMESTAMP WITH TIME ZONE,
		ADD COLUMN is_active BOOLEAN,
		ADD COLUMN is_deleted BOOLEAN,
		ADD COLUMN created_by_id UUID,
		ADD COLUMN updated_by_id UUID,
		ADD COLUMN created_at TIMESTAMP WITH TIME ZONE,
		ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE;`

	_, err = tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// DropUsersTable rollback
func (m *mig) DropAccountsTable() error {
	tx := m.GetTx()

	st := `DROP TABLE accounts;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

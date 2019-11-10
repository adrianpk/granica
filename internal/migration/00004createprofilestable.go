package migration

import "log"

// CreateProfilesTable migration
func (m *mig) CreateProfilesTable() error {
	tx := m.GetTx()

	st := `CREATE TABLE profiles

	(
		id UUID PRIMARY KEY,
		tenant_id VARCHAR(128),
		slug VARCHAR(36) UNIQUE,
		owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
	  account_type VARCHAR(36),
		name VARCHAR(64),
		email VARCHAR(255) UNIQUE,
		description TEXT NULL,
		location VARCHAR(255) NULL,
		bio VARCHAR(255),
		moto VARCHAR(255),
		website VARCHAR(255),
		aniversary_date TIMESTAMP,
		avatar_path VARCHAR(255),
		header_path VARCHAR(255),
		additiional_data jsonb,
	);`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	st = `
		ALTER TABLE profiles
		ADD COLUMN geolocation geography (Point,4326),
		ADD COLUMN locale VARCHAR(32),
		ADD COLUMN base_tz VARCHAR(2),
		ADD COLUMN current_tz VARCHAR(2),
		ADD COLUMN is_active BOOLEAN,
		ADD COLUMN is_deleted BOOLEAN,
		ADD COLUMN created_by_id UUID REFERENCES users(id),
		ADD COLUMN updated_by_id UUID REFERENCES users(id),
		ADD COLUMN created_at TIMESTAMP WITH TIME ZONE,
		ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE;`

	_, err = tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// DropProfilesTable rollback
func (m *mig) DropProfilesTable() error {
	tx := m.GetTx()

	st := `DROP TABLE profiles;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

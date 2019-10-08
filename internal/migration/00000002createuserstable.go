package migration

import "log"

const (
	name2 = "create_users_table"
)

// Up00000002 migration
func (m *mig) Up00000002() (string, error) {
	tx := m.GetTx()

	st := `CREATE TABLE users
	(
		id UUID PRIMARY KEY,
		slug VARCHAR(36),
		username VARCHAR(32) UNIQUE,
		password_digest CHAR(128),
		email VARCHAR(255) UNIQUE,
		given_name VARCHAR(32),
		middle_names VARCHAR(32) NULL,
		family_name VARCHAR(64)
	);`

	_, err := tx.Exec(st)
	if err != nil {
		return name2, err
	}

	st = `
		ALTER TABLE users
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
		return name2, err
	}

	return name2, nil
}

// Down00000002 migration
func (m *mig) Down00000002() (string, error) {
	tx := m.GetTx()

	st := `DROP TABLE users;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return name2, err
	}

	return name2, nil
}

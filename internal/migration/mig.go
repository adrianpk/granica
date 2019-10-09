package migration

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/mikrowezel/backend/migration"
)

type (
	mig struct {
		name string
		up   migration.Fx
		down migration.Fx
		tx   *sqlx.Tx
	}
)

func (m *mig) Config(up migration.Fx, down migration.Fx) {
	m.up = up
	m.down = down
}

func (m *mig) GetName() (name string) {
	return m.name
}

func (m *mig) GetUp() (up migration.Fx) {
	return m.up
}

func (m *mig) GetDown() (down migration.Fx) {
	return m.down
}

func (m *mig) SetTx(tx *sqlx.Tx) {
	m.tx = tx
}

func (m *mig) GetTx() (tx *sqlx.Tx) {
	return m.tx
}

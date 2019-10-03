package migration

import "github.com/jmoiron/sqlx"

type (
	mig struct {
		fx  func() (string, error)
		tx  *sqlx.Tx
		err error
	}
)

func (m *mig) SetFx(fx func() (string, error)) {
	m.fx = fx
}

func (m *mig) SetTx(tx *sqlx.Tx) {
	m.tx = tx
}

func (m *mig) GetTx() *sqlx.Tx {
	return m.tx
}

func (m *mig) SetErr(err error) {
	m.err = err
}

func (m *mig) GetErr() error {
	return m.err
}

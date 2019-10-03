package migration

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/config"
)

type (
	migrator struct {
		cfg  *config.Config
		conn *sqlx.DB
		up   []*migration
		down []*rollback
	}

	proc struct {
		order    int
		tx       transaction
		executed bool
		err      error
	}

	transaction struct {
		conn     *sqlx.DB
		function func() procResult
	}

	migration struct {
		proc
	}

	rollback struct {
		proc
	}

	procResult struct {
		tx   *sqlx.Tx
		name string
		err  error
	}
)

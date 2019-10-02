package migration

import "github.com/jmoiron/sqlx"

type (
	Migrator struct {
		conn *sqlx.DB
		up   []*Migration
		down []*Rollback
	}

	proc struct {
		order    int
		tx       transaction
		executed bool
		err      error
	}

	transaction struct {
		conn     *sqlx.DB
		function func() error
	}

	Migration struct {
		proc
	}

	Rollback struct {
		proc
	}
)

package postgres

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
)

type retryResult struct {
	conn *sqlx.DB
	err  error
}

// RetryConnection implments a backoff mechanism for establishing a connection
// to Postgres; this is especially useful in containerized environments where
// components can be started out of order.
func (h *DbHandler) RetryConnection() chan retryResult {
	res := make(chan retryResult)

	cbmax := uint64(h.Cfg().ValAsInt("pg.backoff.maxtries", 1))
	bo := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), cbmax)

	go func() {
		defer close(res)

		url := h.dbURL()

		for i := 0; i <= int(cbmax); i++ {

			h.Log().Info("Dialing to Postgres", "host", url)

			conn, err := sqlx.Open("postgres", url)
			if err != nil {
				h.Log().Error(err, "Postgres connection error")
			}

			err = conn.Ping()
			if err == nil {
				h.Log().Info("Postgres connection established")
				res <- retryResult{conn, nil}
				return
			}

			h.Log().Error(err, "Postgres connection error")

			// Backoff
			nb := bo.NextBackOff()
			if nb == backoff.Stop {
				h.Log().Info("Postgres connection failed", "reason", "max number of attempts reached")
				err := errors.New("Postgres max number of connection attempts reached")
				res <- retryResult{nil, err}
				bo.Reset()
				return
			}

			h.Log().Info("Postgres connection failed", "retrying-in", nb.String(), "unit", "seconds")
			time.Sleep(nb)
		}
	}()

	return res
}

func (h *DbHandler) dbURL() string {
	host := h.Cfg().ValOrDef("pg.host", "")
	port := h.Cfg().ValAsInt("pg.port", 5432)
	db := h.Cfg().ValOrDef("pg.database", "")
	user := h.Cfg().ValOrDef("pg.user", "")
	pass := h.Cfg().ValOrDef("pg.password", "")
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, db)
}

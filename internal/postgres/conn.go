package postgres

import (
	"time"

	_ "github.com/lib/pq"

	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
)

// RetryConnection implements a backoff mechanism for establishing a connection
// to Postgres; this is especially useful in containerized environments where
// components can be started out of order.
func (h *DbHandler) RetryConnection() chan *sqlx.DB {
	result := make(chan *sqlx.DB)

	cbmax := uint64(h.Cfg().ValAsInt("postgres.backoff.maxtries", 1))
	bo := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), cbmax)

	go func() {
		defer close(result)

		url := h.Cfg().ValOrDef("postgres.conn.url", "")

		for i := 0; i <= int(cbmax); i++ {

			h.Log().Info("Dialing to Postgres", "host", url)

			conn, err := sqlx.Open("postgres", url)
			if err != nil {
				h.Log().Error(err, "Postgres connection error")
			}

			err = conn.Ping()
			if err == nil {
				h.Log().Info("Postgres connection established")
				result <- conn
				return
			}

			h.Log().Error(err, "Postgres connection error")

			// Backoff
			nb := bo.NextBackOff()
			if nb == backoff.Stop {
				result <- nil
				h.Log().Info("Postgres connection failed", "reason", "max number of tries reached")
				bo.Reset()
				return
			}

			h.Log().Info("Postgres connection failed", "retrying-in", nb.String(), "unit", "seconds")
			time.Sleep(nb)
		}
	}()

	return result
}

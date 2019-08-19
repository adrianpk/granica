package cockroach

import (
	"time"

	_ "github.com/lib/pq"

	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
)

// RetryConnection implements a backoff mechanism for establishing a connection
// to Cockroach; this is especially useful in containerized environments where
// components can be started out of ordeh.
func (h *DbHandler) RetryConnection() chan *sqlx.DB {
	result := make(chan *sqlx.DB)

	cbmax := uint64(h.Cfg().ValAsInt("cockroach.backoff.maxtries", 1))
	bo := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), cbmax)

	go func() {
		defer close(result)

		url := h.Cfg().ValOrDef("cockroach.conn.url", "postgres://root@localhost:26257/mydb?sslmode=disable")

		for i := 0; i <= int(cbmax); i++ {

			h.Log().Info("Dialing to Coockroach", "host", url)

			conn, err := sqlx.Open("cockroach", url)
			if err != nil {
				h.Log().Error(err, "Cockroach connection error")
			}

			err = conn.Ping()
			if err == nil {
				h.Log().Info("Cockroach connection established")
				result <- conn
				return
			}

			h.Log().Error(err, "Cockroach connection error")

			// Backoff
			nb := bo.NextBackOff()
			if nb == backoff.Stop {
				result <- nil
				h.Log().Info("Cockroach connection failed", "reason", "max number of tries reached")
				bo.Reset()
				return
			}

			h.Log().Info("Cockroach connection failed", "retrying-in", nb.String(), "unit", "seconds")
			time.Sleep(nb)
		}
	}()

	return result
}

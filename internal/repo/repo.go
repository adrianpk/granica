package repo

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/db/postgres"
	"gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

var (
	// Repo is a package level repo handler instance.
	Handler *Repo
)

type (
	// Repo is a repo handler.
	Repo struct {
		*postgres.DbHandler
		Tx *sqlx.Tx
	}
)

// NewRepo creates and returns a new repo handler.
func NewHandler(ctx context.Context, cfg *config.Config, log *log.Logger, name string) (*Repo, error) {
	if name == "" {
		name = fmt.Sprintf("repo-handler-%s", nameSufix())
	}
	log.Info("New handler", "name", name)

	dbh, err := postgres.NewHandler(ctx, cfg, log, name)
	if err != nil {
		return nil, err
	}

	h := &Repo{
		DbHandler: dbh,
	}

	return h, nil
}

// Init a new repo handler.
// it also stores it as the package default handler.
func (h *Repo) Init(s svc.Service) chan bool {
	// Set package default handler.
	// TODO: See if this could be avoided.
	Handler = h

	ok := make(chan bool)
	go func() {
		defer close(ok)
		conn, err := h.Connect()
		if err != nil {
			s.Log().Error(err, "Init Postgres Db handler error")
			ok <- false
			return
		}
		h.Conn = conn
		s.Lock()
		s.AddHandler(h)
		s.Unlock()
		h.Log().Info("Repo initializated", "name", h.Name())
		ok <- true
	}()
	return ok
}

// GetTx returns repo current transaction.
// Creates a new one if it is nil.
func (r *Repo) GetTx() (*sqlx.Tx, error) {
	if r.Tx == nil {
		return r.InitTx()
	}
	return r.Tx, nil
}

// InitTx initializes a transaction.
func (r *Repo) InitTx() (*sqlx.Tx, error) {
	tx, err := r.Conn.Beginx()
	if err != nil {
		return nil, err
	}

	r.Tx = tx
	return tx, err
}

// CommitTx commits a transaction.
func (r *Repo) CommitTx() error {
	if r.Tx == nil {
		return errors.New("no current transaction")
	}

	return r.Tx.Commit()
}

func nameSufix() string {
	digest := hash(time.Now().String())
	return digest[len(digest)-8:]
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}

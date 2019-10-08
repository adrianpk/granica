package repo

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/db/postgres"
	"gitlab.com/mikrowezel/backend/log"
	svc "gitlab.com/mikrowezel/backend/service"
)

var (
	// Repo is a package level repo handler instance.
	Handler *Repo
)

type (
	// Repo is a repo handler.
	Repo struct {
		*postgres.DbHandler
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
		_, err := h.Connect()
		if err != nil {
			s.Log().Error(err, "Init Postgres Db handler error")
			ok <- false
			return
		}
		s.Lock()
		s.AddHandler(h)
		s.Unlock()
		h.Log().Info("Repo initializated", "name", h.Name())
		ok <- true
	}()
	return ok
}

// NewTx returns a new transcation.
func (r *Repo) NewTx() (*sqlx.Tx, error) {
	return r.Conn.Beginx()
}

func (r *Repo) UserRepo(tx *sqlx.Tx) *UserRepo {
	return makeUserRepo(context.Background(), r.Cfg(), r.Log(), tx)
}

func (r *Repo) UserRepoNewTx() (*UserRepo, error) {
	tx, err := r.NewTx()
	if err != nil {
		return nil, err
	}
	return makeUserRepo(context.Background(), r.Cfg(), r.Log(), tx), nil
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

package repo

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/db/postgres"
	"gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

var (
	// Repo is a package level repo handler instance.
	Repo *RepoHandler
)

type (
	// RepoHandler is a repo handler.
	RepoHandler struct {
		*postgres.DbHandler
	}
)

// NewHandler creates and returns a new repo handler.
func NewHandler(ctx context.Context, cfg *config.Config, log *log.Logger) (*RepoHandler, error) {
	n := fmt.Sprintf("repo-handler-%s", nameSufix())
	log.Info("New handler", "name", n)

	dbh, err := postgres.NewHandler(ctx, cfg, log, n)
	if err != nil {
		return nil, err
	}

	h := &RepoHandler{
		DbHandler: dbh,
	}

	return h, nil
}

// In
// it also stores it as the package default handler.
func (h *RepoHandler) Init(s svc.Service) chan bool {
	h.DbHandler.Init(s)
	// Set package default handler.
	// TODO: See if this could be avoided.
	Repo = h

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
		h.Log().Info("Handler initializated", "name", h.Name())
		ok <- true
	}()
	return ok
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

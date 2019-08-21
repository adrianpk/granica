package postgres

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	_ "github.com/go-sql-driver/mysql" // package init.
	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

var (
	// Db is a package level DB handler instance.
	Db *DbHandler
)

type (
	// DbHandler is a DB handler.
	DbHandler struct {
		*svc.BaseHandler
		Conn *sqlx.DB
	}
)

// NewHandler creates and returns a new DB handler.
func NewHandler(ctx context.Context, cfg *config.Config, log *log.Logger) (*DbHandler, error) {
	hn := fmt.Sprintf("postgres-db-handler-%s", nameSufix())
	log.Info("New handler", "name", hn)
	h := &DbHandler{
		BaseHandler: svc.NewBaseHandler(ctx, cfg, log, hn),
	}
	return h, nil
}

// Init a new Postgres database handler.
// it also stores it as the package default handler.
func (h *DbHandler) Init(s svc.Service) chan bool {
	// Set package default handler.
	// TODO: See if this could be avoided.
	Db = h

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

// Start Postgres database handler.
func (h *DbHandler) Start() error {
	h.Enable()
	return nil
}

// Connect to Posgres database.
func (h *DbHandler) Connect() (*sqlx.DB, error) {
	res := <-h.RetryConnection()
	if res.err != nil {
		return nil, res.err
	}
	return res.conn, nil
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

package cockroach

import (
	"context"

	_ "github.com/go-sql-driver/mysql" // package init.
	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

var (
	// Db is a package general use DB handler.
	Db *DbHandler
)

type (
	// Handler is a DB handler.
	DbHandler struct {
		*svc.BaseHandler
		Conn *sqlx.DB
	}
)

// InitDB creates and return a new DB handler.
// it also stores it as the package default handler.
func InitDB(ctx context.Context, cfg *config.Config, log *log.Logger) (*DbHandler, error) {
	var err error
	Db, err = newHandler(ctx, cfg, log)
	if err != nil {
		return nil, err
	}
	return Db, nil
}

// NewHandler creates and returns a new DB handler.
func newHandler(ctx context.Context, cfg *config.Config, log *log.Logger) (*DbHandler, error) {
	h := &DbHandler{
		BaseHandler: svc.NewBaseHandler(ctx, cfg, log, "cockroach-db-handler"),
	}

	//h.Conn <- h.RetryConnection()

	return h, nil // TODO: RetryConnection will eventually throw a timeout error.
}

//// GetConn from MySQL.
//func GetConn() (*sqlx.DB, error) {
//if pctx == nil || pcfg == nil || plog == nil {
//return nil, errors.New("package has no initialized yet")
//}

//if DB == nil || DB.Conn == nil || DB.Conn.Ping() != nil {
//InitDB(pctx, pcfg, plog)
//}

//return DB.Conn, nil
//}

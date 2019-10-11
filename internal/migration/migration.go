package migration

import (
	"context"
	"fmt"

	_ "github.com/lib/pq" // package init.
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/log"
	"gitlab.com/mikrowezel/backend/migration"

	svc "gitlab.com/mikrowezel/backend/service"
)

const (
	devDb  = "granica"
	testDb = "granica_test"
	prodDb = "granica_prod"
)

type (
	// Migrator is a migrator handler.
	Migrator struct {
		*svc.BaseHandler
		*migration.Migrator
	}
)

// NewHandler creates and returns a new repo handler.
func NewHandler(ctx context.Context, cfg *config.Config, log *log.Logger, name string) (*Migrator, error) {
	if name == "" {
		name = fmt.Sprintf("migration-handler-%s", svc.NameSufix())
	}
	log.Info("New handler", "name", name)

	h := &Migrator{
		BaseHandler: svc.NewBaseHandler(ctx, cfg, log, name),
		Migrator:    GetMigrator(cfg),
	}

	return h, nil
}

// Init a new repo handler.
// it also stores it as the package default handler.
func (h *Migrator) Init(s svc.Service) chan bool {
	// Set package default handler.
	// TODO: See if this could be avoided.
	ok := make(chan bool)
	go func() {
		defer close(ok)

		// NOTE: Remove this before release
		err := h.SoftReset()
		if err != nil {
			s.Log().Error(err, "Init Postgres Db handler error")
			ok <- false
			return
		}

		err = h.Migrate()
		if err != nil {
			s.Log().Error(err, "Init Postgres Db handler error")
			ok <- false
			return
		}

		s.Lock()
		s.AddHandler(h)
		s.Unlock()
		h.Log().Info("Migrator initializated", "name", h.Name())
		ok <- true
	}()
	return ok
}

// GetMigrator configured.
func GetMigrator(cfg *config.Config) *migration.Migrator {
	m := migration.Init(cfg)

	// Migrations
	// Enable Postgis
	mg := &mig{}
	mg.Config(mg.EnablePostgis, mg.DropPostgis)
	m.AddMigration(mg)

	// CreateUsersTable
	mg = &mig{}
	mg.Config(mg.CreateUsersTable, mg.DropUsersTable)
	m.AddMigration(mg)

	return m
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/granica/internal/postgres"
	"gitlab.com/mikrowezel/granica/pkg/auth"
	"gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
	// "gitlab.com/mikrowezel/internal/cockroach"
)

type contextKey string

var (
	s svc.Service
)

func main() {
	cfg := config.Load("grn")
	log := initLog(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go checkSigTerm(cancel)

	// Create service
	s = newService(ctx, cfg, log, cancel)

	// Add service handlers
	db, err := postgres.InitDb(ctx, cfg, log)
	if err != nil {
		log.Error(err, "Cannot create Postgres Db handler")
	}
	s.AddHandler(db)

	// Set service worker
	auth := auth.NewWorker(ctx, cfg, log, "auth-worker")
	s.SetWorker(auth)

	// Initialize handlers and workers
	err = s.Init()
	if err != nil {
		log.Error(err)
		cancel()
	}

	// Start service
	s.Start()

	log.Error(err, "Service stoped")
}

// newService creates a service instance.
func newService(ctx context.Context, cfg *config.Config, log *log.Logger, cancel context.CancelFunc) svc.Service {
	sn := cfg.ValOrDef("svc.name", "granica")
	sr := cfg.ValOrDef("svc.revision", "n/a")
	s := svc.NewService(ctx, cfg, log, cancel, sn, sr)
	return s
}

func initPostgres(s svc.Service) chan bool {
	ok := make(chan bool)
	go func() {
		defer close(ok)
		r, err := postgres.InitDb(s.Ctx(), s.Cfg(), s.Log())
		if err != nil {
			s.Log().Error(err, "Init Postgres Db handler error")
			ok <- false
			return
		}
		s.Lock()
		s.AddHandler(r)
		s.Unlock()
		ok <- true
	}()
	return ok
}

func initLog(cfg *config.Config) *log.Logger {
	ll := int(cfg.ValAsInt("log.level", 1))
	sn := cfg.ValOrDef("svc.name", "granica")
	sr := cfg.ValOrDef("svc.revision", "n/a")
	return log.NewLogger(ll, sn, sr)
}

// checkSigTerm - Listens to sigterm events.
func checkSigTerm(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	s.Stop()
	cancel()
	os.Exit(0)
}

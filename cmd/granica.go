package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/granica/internal/migration"
	"gitlab.com/mikrowezel/backend/granica/internal/repo"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth"
	"gitlab.com/mikrowezel/backend/log"
	svc "gitlab.com/mikrowezel/backend/service"
)

type contextKey string

var (
	s svc.Service
)

func main() {
	cfg := config.Load("grn")
	log := initLog(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	checkStopEvents(ctx, cancel)

	// Create service
	s = newService(ctx, cfg, log, cancel)

	// Add Migration handler
	mh, err := migration.NewHandler(ctx, cfg, log, "migration-handler")
	s.AddHandler(mh)

	// Add Repo handler
	rh, err := repo.NewHandler(ctx, cfg, log, "repo-handler")
	s.AddHandler(rh)

	// Set service worker
	auth, err := auth.NewWorker(ctx, cfg, log, "auth-worker")
	if err != nil {
		exit(log, err)
	}
	s.SetWorker(auth)

	// Initialize handlers and workers
	err = s.Init()
	if err != nil {
		exit(log, err)
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

func initLog(cfg *config.Config) *log.Logger {
	ll := int(cfg.ValAsInt("log.level", 1))
	sn := cfg.ValOrDef("svc.name", "granica")
	sr := cfg.ValOrDef("svc.revision", "n/a")
	//return log.NewLogger(ll, sn, sr)
	return log.NewDevLogger(ll, sn, sr)
}

func exit(log *log.Logger, err error) {
	log.Error(err)
	os.Exit(1)
}

func checkStopEvents(ctx context.Context, cancel context.CancelFunc) {
	go checkSigterm(cancel)
	go checkCancel(ctx)
}

func checkSigterm(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cancel()
}

func checkCancel(ctx context.Context) {
	<-ctx.Done()
	s.Stop()
	os.Exit(1)
}

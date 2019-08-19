package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/mikrowezel/config"
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
	cfg := config.Load("granica")
	log := initLog(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go checkSigTerm(cancel, log)

	s = newService(ctx, cfg, log, cancel)
	w := auth.NewWorker(ctx, cfg, log, "auth-worker")
	w.AttachTo(s)

	err := s.Init()
	if err != nil {
		log.Error(err)
		cancel()
	}

	s.Start()

	log.Error(err, "Service stoped")
}

// newService creates a service instance.
func newService(ctx context.Context, cfg *config.Config, log *log.Logger, cancel context.CancelFunc) svc.Service {
	sn := cfg.ValOrDef("service.name", "granica")
	sv := cfg.ValOrDef("service.version", "n/a")
	s := svc.NewService(ctx, cfg, log, cancel, sn, sv)
	return s
}

func initLog(cfg *config.Config) *log.Logger {
	ll := int(cfg.ValAsInt("log.level", 3))
	sn := cfg.ValOrDef("service.name", "granica")
	sv := cfg.ValOrDef("service.version", "n/a")
	return log.NewLogger(ll, sn, sv)
}

// checkSigTerm - Listens to sigterm events.
func checkSigTerm(cancel context.CancelFunc, log *log.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	s.Stop()
	cancel()
	os.Exit(0)
}

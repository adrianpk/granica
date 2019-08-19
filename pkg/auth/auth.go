package auth

import (
	"context"
	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

type Auth struct {
	*svc.BaseWorker
}

// NewWorker creates a new base worker instance.
// This is a bare implementtion of Worker interface
// just for mocking and/or testing purposes.
func NewWorker(ctx context.Context, cfg *config.Config, log *log.Logger, name string) *Auth {
	w := &Auth{
		BaseWorker: svc.NewWorker(ctx, cfg, log, "granica-auth-worker"),
	}
	return w
}

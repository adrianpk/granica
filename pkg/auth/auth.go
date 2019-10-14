package auth

import (
	"context"
	"net/http"

	"gitlab.com/mikrowezel/backend/config"
	logger "gitlab.com/mikrowezel/backend/log"
	svc "gitlab.com/mikrowezel/backend/service"
)

type (
	Auth struct {
		*svc.BaseWorker
		Server http.Handler
	}

	AuthCtx struct{}

	contextKey string
)

const (
	userCtxKey contextKey = "user"
)

// NewWorker creates a new Auth worker instance.
func NewWorker(ctx context.Context, cfg *config.Config, log *logger.Logger, name string) *Auth {
	w := &Auth{
		BaseWorker: svc.NewWorker(ctx, cfg, log, "granica-auth-worker"),
	}
	w.AddServer()
	return w
}

func (a *Auth) Start() error {
	p := a.Cfg().ValOrDef("server.port", ":8080")
	a.Log().Info("Server initializing", "port", p)
	err := http.ListenAndServe(p, a.Server)
	a.Log().Error(err)
	return err
}

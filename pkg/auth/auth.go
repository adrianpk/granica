package auth

import (
	"context"
	"net/http"

	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/jsonrest"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
	logger "gitlab.com/mikrowezel/backend/log"
	svc "gitlab.com/mikrowezel/backend/service"
)

type (
	Auth struct {
		*svc.BaseWorker
		service *service.Service
		jsonep  *jsonrest.Endpoint
		Server  http.Handler
	}
)

// NewWorker creates a new Auth worker instance.
func NewWorker(ctx context.Context, cfg *config.Config, log *logger.Logger, name string) *Auth {

	service := service.MakeService(ctx, cfg, log)

	w := &Auth{
		BaseWorker: svc.NewWorker(ctx, cfg, log, "granica-auth-worker"),
		service:    service,
		jsonep:     jsonrest.MakeEndpoint(ctx, cfg, log, service),
	}

	w.AddServer()
	return w
}

func (a *Auth) Init() bool {
	rh, err := a.repoHandler()
	if err != nil {
		return false
	}
	a.service.SetRepo(rh)
	return true
}

func (a *Auth) Start() error {
	p := a.Cfg().ValOrDef("server.port", ":8080")
	a.Log().Info("Server initializing", "port", p)
	err := http.ListenAndServe(p, a.Server)
	a.Log().Error(err)
	return err
}

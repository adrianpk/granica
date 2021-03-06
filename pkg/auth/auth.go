package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/jsonrest"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/web"
	logger "gitlab.com/mikrowezel/backend/log"
	svc "gitlab.com/mikrowezel/backend/service"
)

type (
	Auth struct {
		*svc.BaseWorker
		service        *service.Service
		i18nBundle     *i18n.Bundle
		webep          *web.Endpoint
		jsonep         *jsonrest.Endpoint
		WebServer      http.Handler
		JSONRESTServer http.Handler
	}
)

// NewWorker creates a new Auth worker instance.
func NewWorker(ctx context.Context, cfg *config.Config, log *logger.Logger, name string) (*Auth, error) {
	service := service.MakeService(ctx, cfg, log)

	wep, err := web.MakeEndpoint(ctx, cfg, log, service)
	if err != nil {
		return nil, err
	}

	w := &Auth{
		BaseWorker: svc.NewWorker(ctx, cfg, log, "granica-auth-worker"),
		service:    service,
		webep:      wep,
		jsonep:     jsonrest.MakeEndpoint(ctx, cfg, log, service),
	}

	w.AddWebServer()
	w.AddJSONRESTServer()

	return w, nil
}

func (a *Auth) Init() bool {
	rh, err := a.repoHandler()
	if err != nil {
		return false
	}
	a.service.SetRepo(rh)

	mlh, err := a.mailerHandler()
	if err != nil {
		return false
	}
	a.service.SetMailer(mlh)
	return true
}

func (a *Auth) Start() error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		a.StartWeb()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		a.StartJSONREST()
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func (a *Auth) StartWeb() error {
	p := a.Cfg().ValOrDef("web.server.port", "8080")
	p = fmt.Sprintf(":%s", p)

	a.Log().Info("Web server initializing", "port", p)

	err := http.ListenAndServe(p, a.WebServer)
	a.Log().Error(err)
	return err
}

func (a *Auth) StartJSONREST() error {
	p := a.Cfg().ValOrDef("jsonrest.server.port", "8081")
	p = fmt.Sprintf(":%s", p)

	a.Log().Info("JSON REST Server initializing", "port", p)

	err := http.ListenAndServe(p, a.JSONRESTServer)
	a.Log().Error(err)
	return err
}

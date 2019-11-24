package web

import (
	"context"
	"encoding/gob"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/mikrowezel/backend/config"
	svc "gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
	"gitlab.com/mikrowezel/backend/log"
	"gitlab.com/mikrowezel/backend/web"
)

type (
	Endpoint struct {
		*web.Endpoint
		service *svc.Service
		i18n    *i18n.Bundle
	}
)

func MakeEndpoint(ctx context.Context, cfg *config.Config, log *log.Logger, s *svc.Service) (*Endpoint, error) {
	registerGobTypes()

	wep, err := web.MakeEndpoint(ctx, cfg, log, pathFxs)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		Endpoint: wep,
		service:  s,
	}, nil
}

func registerGobTypes() {
	gob.Register(tp.User{})
}

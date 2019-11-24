package web

import (
	"context"

	"github.com/adriank/go-i18n/v2/i18n"
	"gitlab.com/mikrowezel/backend/config"
	svc "gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
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
	wep, err := web.MakeEndpoint(ctx, cfg, log, pathFxs)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		Endpoint: wep,
		service:  s,
	}, nil
}

// TODO: Compare approachs: passing i18n localizer as argument vs. get it from request context through middleware.
//func MakeEndpoint(ctx context.Context, cfg *config.Config, log *log.Logger, s *svc.Service, b *i18n.Bundle) (*Endpoint, error) {
//wep, err := web.MakeEndpoint(ctx, cfg, log, pathFxs)
//if err != nil {
//return nil, err
//}

//return &Endpoint{
//Endpoint: wep,
//service:  s,
//i18n:     b,
//}, nil
//}

package web

import (
	"context"

	"gitlab.com/mikrowezel/backend/config"
	svc "gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
	"gitlab.com/mikrowezel/backend/log"
	"gitlab.com/mikrowezel/backend/web"
)

type (
	Endpoint struct {
		*web.Endpoint
		service *svc.Service
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

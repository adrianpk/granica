package web

import (
	"context"

	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
	"gitlab.com/mikrowezel/backend/log"
	"gitlab.com/mikrowezel/backend/web"
)

type (
	Endpoint struct {
		*web.Endpoint
	}
)

func MakeEndpoint(ctx context.Context, cfg *config.Config, log *log.Logger, service *service.Service) (*Endpoint, error) {
	wep, err := web.MakeEndpoint(ctx, cfg, log, service, pathFxs)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		wep,
	}, nil
}

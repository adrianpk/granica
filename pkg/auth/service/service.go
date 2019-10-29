package service

import (
	"context"

	"gitlab.com/mikrowezel/backend/log"

	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/granica/internal/repo"
)

type Service struct {
	ctx  context.Context
	cfg  *config.Config
	log  *log.Logger
	repo *repo.Repo
}

func MakeService(ctx context.Context, cfg *config.Config, log *log.Logger) *Service {
	return &Service{
		ctx: ctx,
		cfg: cfg,
		log: log,
	}
}

func (s *Service) Ctx() context.Context {
	return s.ctx
}

func (s *Service) Cfg() *config.Config {
	return s.cfg
}

func (s *Service) Log() *log.Logger {
	return s.log
}

// Repo
func (s *Service) SetRepo(repo *repo.Repo) {
	s.repo = repo
}

func (s *Service) userRepo() (*repo.UserRepo, error) {
	return s.repo.UserRepoNewTx()
}

func (s *Service) accountRepo() (*repo.AccountRepo, error) {
	return s.repo.AccountRepoNewTx()
}

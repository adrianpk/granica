package auth

import (
	"errors"

	"gitlab.com/mikrowezel/backend/granica/internal/mailer"
	"gitlab.com/mikrowezel/backend/granica/internal/repo"
)

// TODO: Move functions to a more appropriate place.

// Repo
func (a *Auth) repoHandler() (*repo.Repo, error) {
	h, ok := a.Handler("repo-handler")
	if !ok {
		return nil, errors.New("repo handler not available")
	}

	repo, ok := h.(*repo.Repo)
	if !ok {
		return nil, errors.New("invalid repo handler")
	}

	return repo, nil
}

func (a *Auth) userRepo() (*repo.UserRepo, error) {
	rh, err := a.repoHandler()
	if err != nil {
		return nil, err
	}
	return rh.UserRepoNewTx()
}

func (a *Auth) accountRepo() (*repo.AccountRepo, error) {
	rh, err := a.repoHandler()
	if err != nil {
		return nil, err
	}
	return rh.AccountRepoNewTx()
}


// Mailer
func (a *Auth) mailerHandler() (*mailer.SESMailer, error) {
	h, ok := a.Handler("mailer-handler")
	if !ok {
		return nil, errors.New("mailer handler not available")
	}

	mailer, ok := h.(*mailer.SESMailer)
	if !ok {
		return nil, errors.New("invalid mailer handler")
	}

	return mailer, nil
}

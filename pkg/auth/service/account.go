package service

import (
	"gitlab.com/mikrowezel/backend/granica/internal/repo"
	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
)

const (
	createAccountErr = "cannot create account"
	getAllAccountErr = "cannot get account list"
	getAccountErr    = "cannot get account"
	updateAccountErr = "cannot update account"
	deleteAccountErr = "cannot delete account"
)

func (s *Service) CreateAccount(req tp.CreateAccountReq, res *tp.CreateAccountRes) error {
	// Model
	u := req.ToModel()

	// Repo
	repo, err := s.accountRepo()
	if err != nil {
		res.FromModel(nil, createAccountErr, err)
		return err
	}

	err = repo.Create(&u)
	if err != nil {
		res.FromModel(nil, createAccountErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, createAccountErr, err)
		return err
	}

	// Output
	res.FromModel(&u, "", nil)
	return nil
}

func (s *Service) GetAccounts(req tp.GetAccountsReq, res *tp.GetAccountsRes) error {
	// Repo
	repo, err := s.accountRepo()
	if err != nil {
		res.FromModel(nil, getAllAccountErr, err)
		return err
	}

	us, err := repo.GetAll()
	if err != nil {
		res.FromModel(nil, getAllAccountErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, getAllAccountErr, err)
		return err
	}

	// Output
	res.FromModel(us, "", nil)
	return nil
}

func (s *Service) GetAccount(req tp.GetAccountReq, res *tp.GetAccountRes) error {
	// Model
	u := req.ToModel()

	// Repo
	repo, err := s.accountRepo()
	if err != nil {
		res.FromModel(nil, getAccountErr, err)
		return err
	}

	u, err = repo.GetBySlug(u.Slug.String)
	if err != nil {
		res.FromModel(nil, getAccountErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, getAccountErr, err)
		return err
	}

	// Output
	res.FromModel(&u, "", nil)
	return nil
}

func (s *Service) UpdateAccount(req tp.UpdateAccountReq, res *tp.UpdateAccountRes) error {
	// Repo
	repo, err := s.accountRepo()
	if err != nil {
		res.FromModel(nil, updateAccountErr, err)
		return err
	}

	// Get account
	current, err := repo.GetBySlug(req.Identifier.Slug)
	if err != nil {
		res.FromModel(nil, updateAccountErr, err)
		return err
	}

	// Create a model
	u := req.ToModel()
	u.ID = current.ID

	// Update
	err = repo.Update(&u)
	if err != nil {
		res.FromModel(nil, updateAccountErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, updateAccountErr, err)
		return err
	}

	// Output
	res.FromModel(&u, "", nil)
	return nil
}

func (s *Service) DeleteAccount(req tp.DeleteAccountReq, res *tp.DeleteAccountRes) error {
	// Repo
	repo, err := s.accountRepo()
	if err != nil {
		res.FromModel(nil, updateAccountErr, err)
		return err
	}

	err = repo.DeleteBySlug(req.Slug)
	if err != nil {
		res.FromModel(nil, updateAccountErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, updateAccountErr, err)
		return err
	}

	// Output
	res.FromModel(nil, "", nil)
	return nil
}

// Misc
func (s *Service) accountRepo() (*repo.AccountRepo, error) {
	return s.repo.AccountRepoNewTx()
}

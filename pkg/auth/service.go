package auth

import (
	"errors"

	"gitlab.com/mikrowezel/granica/internal/repo"
)

func (a *Auth) createUser(req CreateUserReq, res *CreateUserRes) error {
	// Model
	u := req.toModel()

	// Repo
	repo, err := a.userRepo()
	if err != nil {
		res.fromModel(nil, createErr, err)
		return err
	}

	err = repo.Create(&u)
	if err != nil {
		res.fromModel(nil, createErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.fromModel(nil, createErr, err)
		return err
	}

	// Output
	res.fromModel(&u, "", nil)
	return nil
}

func (a *Auth) getUsers(req GetUsersReq, res *GetUsersRes) error {
	// Repo
	repo, err := a.userRepo()
	if err != nil {
		res.fromModel(nil, getAllErr, err)
		return err
	}

	us, err := repo.GetAll()
	if err != nil {
		res.fromModel(nil, getAllErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.fromModel(nil, getAllErr, err)
		return err
	}

	// Output
	res.fromModel(us, "", nil)
	return nil
}

func (a *Auth) getUser(req GetUserReq, res *GetUserRes) error {
	// Model
	u := req.toModel()

	// Repo
	repo, err := a.userRepo()
	if err != nil {
		res.fromModel(nil, getErr, err)
		return err
	}

	u, err = repo.GetByUsername(u.Username.String)
	if err != nil {
		res.fromModel(nil, getErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.fromModel(nil, getErr, err)
		return err
	}

	// Output
	res.fromModel(&u, "", nil)
	return nil
}

func (a *Auth) updateUser(req UpdateUserReq, res *UpdateUserRes) error {
	// Repo
	repo, err := a.userRepo()
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	// Get user
	current, err := repo.GetByUsername(req.Identifier.Username)
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	// Create a model
	// Neither ID nor Username should change.
	u := req.toModel()
	u.ID = current.ID
	// Set envar GRN_APP_USERNAME_UPDATABLE=true
	// to let username be updatable.
	if !(a.Cfg().ValAsBool("app.username.updatable", false)) {
		u.Username = current.Username
	}

	// Update
	err = repo.Update(&u)
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	// Output
	res.fromModel(&u, "", nil)
	return nil
}

func (a *Auth) deleteUser(req DeleteUserReq, res *DeleteUserRes) error {
	// Repo
	repo, err := a.userRepo()
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	err = repo.DeleteByUsername(req.Identifier.Username)
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	// Output
	res.fromModel(nil, "", nil)
	return nil
}

func (a *Auth) userRepo() (*repo.UserRepo, error) {
	rh, err := a.repoHandler()
	if err != nil {
		return nil, err
	}
	return rh.UserRepoNewTx()
}

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

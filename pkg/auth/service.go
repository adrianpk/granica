package auth

// User -----------------------------------------------------------------------
func (a *Auth) CreateUser(req CreateUserReq, res *CreateUserRes) error {
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

func (a *Auth) GetUsers(req GetUsersReq, res *GetUsersRes) error {
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

func (a *Auth) GetUser(req GetUserReq, res *GetUserRes) error {
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

func (a *Auth) UpdateUser(req UpdateUserReq, res *UpdateUserRes) error {
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

func (a *Auth) DeleteUser(req DeleteUserReq, res *DeleteUserRes) error {
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

// Account  -------------------------------------------------------------------
func (a *Auth) CreateAccount(req CreateAccountReq, res *CreateAccountRes) error {
	// Model
	u := req.toModel()

	// Repo
	repo, err := a.accountRepo()
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

func (a *Auth) GetAccounts(req GetAccountsReq, res *GetAccountsRes) error {
	// Repo
	repo, err := a.accountRepo()
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

func (a *Auth) GetAccount(req GetAccountReq, res *GetAccountRes) error {
	// Model
	u := req.toModel()

	// Repo
	repo, err := a.accountRepo()
	if err != nil {
		res.fromModel(nil, getErr, err)
		return err
	}

	u, err = repo.GetByAccountname(u.Slug.String)
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

func (a *Auth) UpdateAccount(req UpdateAccountReq, res *UpdateAccountRes) error {
	// Repo
	repo, err := a.accountRepo()
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	// Get account
	current, err := repo.GetByAccountname(req.Identifier.Slug)
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	// Create a model
	u := req.toModel()
	u.ID = current.ID

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

func (a *Auth) DeleteAccount(req DeleteAccountReq, res *DeleteAccountRes) error {
	// Repo
	repo, err := a.accountRepo()
	if err != nil {
		res.fromModel(nil, updateErr, err)
		return err
	}

	err = repo.DeleteBySlug(req.Identifier.Slug)
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

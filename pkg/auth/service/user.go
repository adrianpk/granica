package service

import (
	"gitlab.com/mikrowezel/backend/granica/internal/repo"
	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
)

const (
	createUserErr = "cannot create resource"
	getAllUserErr = "cannot get resource list"
	getUserErr    = "cannot get resource"
	updateUserErr = "cannot update resource"
	deleteUserErr = "cannot delete resource"
)

func (s *Service) CreateUser(req tp.CreateUserReq, res *tp.CreateUserRes) error {
	// Model
	u := req.ToModel()

	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, createUserErr, err)
		return err
	}

	err = repo.Create(&u)
	if err != nil {
		res.FromModel(nil, createUserErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, createUserErr, err)
		return err
	}

	// Output
	res.FromModel(&u, "", nil)
	return nil
}

func (s *Service) IndexUsers(req tp.IndexUsersReq, res *tp.IndexUsersRes) error {
	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, getAllUserErr, err)
		return err
	}

	us, err := repo.GetAll()
	if err != nil {
		res.FromModel(nil, getAllUserErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, getAllUserErr, err)
		return err
	}

	// Output
	res.FromModel(us, "", nil)
	return nil
}

func (s *Service) GetUser(req tp.GetUserReq, res *tp.GetUserRes) error {
	// Model
	u := req.ToModel()

	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, getUserErr, err)
		return err
	}

	u, err = repo.GetBySlug(u.Slug.String)
	if err != nil {
		res.FromModel(nil, getUserErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, getUserErr, err)
		return err
	}

	// Output
	res.FromModel(&u, "", nil)
	return nil
}

func (s *Service) GetUserByUsername(req tp.GetUserReq, res *tp.GetUserRes) error {
	// Model
	u := req.ToModel()

	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, getUserErr, err)
		return err
	}

	u, err = repo.GetByUsername(u.Username.String)
	if err != nil {
		res.FromModel(nil, getUserErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, getUserErr, err)
		return err
	}

	// Output
	res.FromModel(&u, "", nil)
	return nil
}

func (s *Service) UpdateUser(req tp.UpdateUserReq, res *tp.UpdateUserRes) error {
	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, updateUserErr, err)
		return err
	}

	// Get user
	current, err := repo.GetBySlug(req.Identifier.Slug)
	if err != nil {
		res.FromModel(nil, updateUserErr, err)
		return err
	}

	// Create a model
	// Neither ID nor Username should change.
	u := req.ToModel()
	u.ID = current.ID
	// Set envar GRN_APP_USERNAME_UPDATABLE=true
	// to let username be updatable.
	if !(s.Cfg().ValAsBool("app.username.updatable", false)) {
		u.Username = current.Username
	}

	// Update
	err = repo.Update(&u)
	if err != nil {
		res.FromModel(nil, updateUserErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, updateUserErr, err)
		return err
	}

	// Output
	res.FromModel(&u, "", nil)
	return nil
}

func (s *Service) DeleteUser(req tp.DeleteUserReq, res *tp.DeleteUserRes) error {
	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, updateUserErr, err)
		return err
	}

	err = repo.DeleteByUsername(req.Slug)
	if err != nil {
		res.FromModel(nil, updateUserErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(nil, updateUserErr, err)
		return err
	}

	// Output
	res.FromModel(nil, "", nil)
	return nil
}

// Misc
func (s *Service) userRepo() (*repo.UserRepo, error) {
	return s.repo.UserRepoNewTx()
}

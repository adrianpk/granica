package service

import (
	"errors"

	"gitlab.com/mikrowezel/backend/granica/internal/repo"
	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
)

const (
	// Info
	okResultInfo      = "ok_result"
	userCreatedInfo   = "user_created_info"
	userUpdatedInfo   = "user_updated_info"
	userDeletedInfo   = "user_deleted_info"
	userConfirmedInfo = "user_confirmed_info"
	// Error
	cannotProcErr       = "cannot_process_err"
	createUserErr       = "cannot_create_user_err"
	getAllUserErr       = "cannot_get_users_list_err"
	getUserErr          = "cannot_get_user_err"
	updateUserErr       = "cannot_update_user_err"
	deleteUserErr       = "cannot_delete_user_err"
	validationErr       = "validation_error_err"
	signupErr           = "cannot_sign_up_user_err"
	confirmationErr     = "cannot_confirm_user_err"
	signinErr           = "cannot_sign_in_user_err"
	alreadyConfirmedErr = "already_confirm_user_err"
)

func (s *Service) CreateUser(req tp.CreateUserReq, res *tp.CreateUserRes) error {
	// Model
	u := req.ToModel()

	// Validation
	v := NewUserValidator(u)

	err := v.ValidateForCreate()
	if err != nil {
		res.FromModel(&u, validationErr, err)
		return err
	}

	// Confirmation
	u.GenAutoConfirmationToken()

	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(&u, cannotProcErr, err)
		return err
	}

	err = repo.Create(&u)
	if err != nil {
		res.FromModel(&u, createUserErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(&u, createUserErr, err)
		return err
	}

	// Output
	res.FromModel(&u, userCreatedInfo, nil)
	return nil
}

func (s *Service) IndexUsers(req tp.IndexUsersReq, res *tp.IndexUsersRes) error {
	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, cannotProcErr, err)
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
	res.FromModel(us, okResultInfo, nil)
	return nil
}

func (s *Service) GetUser(req tp.GetUserReq, res *tp.GetUserRes) error {
	// Model
	u := req.ToModel()

	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, cannotProcErr, err)
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
	res.FromModel(&u, okResultInfo, nil)
	return nil
}

func (s *Service) GetUserByUsername(req tp.GetUserReq, res *tp.GetUserRes) error {
	// Model
	u := req.ToModel()

	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, cannotProcErr, err)
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
	res.FromModel(&u, okResultInfo, nil)
	return nil
}

func (s *Service) UpdateUser(req tp.UpdateUserReq, res *tp.UpdateUserRes) error {
	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, cannotProcErr, err)
		return err
	}

	// Get user
	current, err := repo.GetBySlug(req.Identifier.Slug)
	if err != nil {
		res.FromModel(nil, getUserErr, err)
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

	// Validation
	v := NewUserValidator(u)

	err = v.ValidateForUpdate()
	if err != nil {
		res.FromModel(nil, validationErr, err)
		return err
	}

	// Update
	err = repo.Update(&u)
	if err != nil {
		res.FromModel(&u, updateUserErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(&u, updateUserErr, err)
		return err
	}

	// Output
	res.FromModel(&u, okResultInfo, nil)
	return nil
}

func (s *Service) DeleteUser(req tp.DeleteUserReq, res *tp.DeleteUserRes) error {
	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(cannotProcErr, err)
		return err
	}

	err = repo.DeleteBySlug(req.Slug)
	if err != nil {
		res.FromModel(deleteUserErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(deleteUserErr, err)
		return err
	}

	// Output
	res.FromModel(okResultInfo, nil)
	return nil
}

func (s *Service) SignUpUser(req tp.SignUpUserReq, res *tp.SignUpUserRes) error {
	// Model
	u := req.ToModel()

	// Validation
	v := NewUserValidator(u)

	err := v.ValidateForSignUp()
	if err != nil {
		res.FromModel(&u, validationErr, err)
	}

	// Generate confirmation token
	u.GenConfirmationToken()

	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(&u, cannotProcErr, err)
		return err
	}

	err = repo.Create(&u)
	if err != nil {
		res.FromModel(&u, cannotProcErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(&u, createUserErr, err)
		return err
	}

	// Mail confirmation
	s.sendConfirmationEmail(&u)

	// Output
	res.FromModel(&u, okResultInfo, nil)
	return nil
}

func (s *Service) ConfirmUser(req tp.GetUserReq, res *tp.GetUserRes) error {
	// Model
	u := req.ToModel()

	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, cannotProcErr, err)
		return err
	}

	s.Log().Debug("Values", "slug", u.Slug.String, "token", u.ConfirmationToken.String)

	u, err = repo.GetBySlugAndToken(u.Slug.String, u.ConfirmationToken.String)
	if err != nil {
		res.FromModel(&u, confirmationErr, err)
		return err
	}

	if u.IsConfirmed.Bool {
		res.FromModel(&u, alreadyConfirmedErr, err)
		return errors.New("already confirmed")
	}

	u, err = repo.ConfirmUser(u.Slug.String, u.ConfirmationToken.String)
	if err != nil {
		res.FromModel(&u, confirmationErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(&u, confirmationErr, err)
		return err
	}

	// Output
	res.FromModel(&u, okResultInfo, nil)
	return nil
}

func (s *Service) SignInUser(req tp.SignInUserReq, res *tp.SignInUserRes) error {
	// Model
	u := req.ToModel()

	// Repo
	repo, err := s.userRepo()
	if err != nil {
		res.FromModel(nil, cannotProcErr, err)
		return err
	}

	u, err = repo.SignIn(u.Username.String, u.Password)
	if err != nil {
		res.FromModel(&u, signinErr, err)
		return err
	}

	err = repo.Commit()
	if err != nil {
		res.FromModel(&u, signinErr, err)
		return err
	}

	// Output
	res.FromModel(&u, okResultInfo, nil)
	return nil
}

// Misc
func (s *Service) userRepo() (*repo.UserRepo, error) {
	return s.repo.UserRepoNewTx()
}

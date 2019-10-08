package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/granica/internal/model"
	logger "gitlab.com/mikrowezel/backend/log"
)

type (
	UserRepo struct {
		ctx context.Context
		cfg *config.Config
		log *logger.Logger
		Tx  *sqlx.Tx
	}
)

func makeUserRepo(ctx context.Context, cfg *config.Config, log *logger.Logger, tx *sqlx.Tx) *UserRepo {
	return &UserRepo{
		ctx: ctx,
		cfg: cfg,
		log: log,
		Tx:  tx,
	}
}

// Create a user in repo.
func (ur *UserRepo) Create(user *model.User) error {
	user.GenID()

	st := `INSERT INTO users (id, slug, username, password_digest, email, given_name, middle_names, family_name, geolocation, locale, base_tz, current_tz, starts_at, ends_at, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
	VALUES (:id, :slug, :username, :password_digest, :email, :given_name, :middle_names, :family_name, :geolocation, :locale, :base_tz, :current_tz, :starts_at, :ends_at, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	_, err := ur.Tx.NamedExec(st, user)

	return err
}

// GetAll users from repo.
func (ur *UserRepo) GetAll() ([]*model.User, error) {
	return nil, nil
}

// Get user data from repo.
func (ur *UserRepo) Get(id interface{}) (*model.User, error) {
	return &model.User{}, nil
}

// Update user data in repo.
func (ur *UserRepo) Update(*model.User) (*model.User, error) {
	return &model.User{}, nil
}

// Delete user data from repo.
func (ur *UserRepo) Delete(id interface{}) error {
	return nil
}

// Commit transaction
func (ur *UserRepo) Commit() error {
	return ur.Tx.Commit()
}

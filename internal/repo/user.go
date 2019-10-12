package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/backend/config"
	logger "gitlab.com/mikrowezel/backend/log"
	"gitlab.com/mikrowezel/granica/internal/model"
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
	user.SetCreateValues()

	st := `INSERT INTO users (id, slug, username, password_digest, email, given_name, middle_names, family_name, geolocation, locale, base_tz, current_tz, starts_at, ends_at, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
	VALUES (:id, :slug, :username, :password_digest, :email, :given_name, :middle_names, :family_name, :geolocation, :locale, :base_tz, :current_tz, :starts_at, :ends_at, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	_, err := ur.Tx.NamedExec(st, user)

	return err
}

// GetAll users from repo.
func (ur *UserRepo) GetAll() (users []model.User, err error) {
	st := `SELECT * FROM public.users;`

	err = ur.Tx.Select(&users, st)

	return users, err
}

// Get user by ID.
func (ur *UserRepo) Get(id interface{}) (model.User, error) {
	var user model.User

	st := `SELECT * FROM USERS WHERE id = '%s' LIMIT 1;`
	st = fmt.Sprintf(st, id.(string))

	err := ur.Tx.Get(&user, st)

	return user, err
}

// GetBySlug user from repo by slug.
func (ur *UserRepo) GetBySlug(slug string) (model.User, error) {
	var user model.User

	st := `SELECT * FROM USERS WHERE slug = '%s' LIMIT 1;`
	st = fmt.Sprintf(st, slug)

	err := ur.Tx.Get(&user, st)

	return user, err
}

// GetByUsername user from repo by username.
func (ur *UserRepo) GetByUsername(username string) (model.User, error) {
	var user model.User

	st := `SELECT * FROM USERS WHERE username = '%s' LIMIT 1;`
	st = fmt.Sprintf(st, username)

	err := ur.Tx.Get(&user, st)

	return user, err
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

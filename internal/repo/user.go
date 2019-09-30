package repo

import (
	"gitlab.com/mikrowezel/granica/internal/model"
)

// CreateUser in repo.
func (r *Repo) CreateUser(user *model.User) error {
	tx, err := r.GetTx()
	if err != nil {
		return err
	}
	st := `INSERT INTO users (id, slug, username, password_digest, email, given_name, middle_names, family_name, geolocation, locale, base_tz, current_tz, starts_at, ends_at, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
	VALUES (:id, :slug, :username, :password_digest, :email, :given_name, :middle_names, :family_name, :geolocation, :locale, :base_tz, :current_tz, :starts_at, :ends_at, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	_, err = tx.NamedExec(st, user)

	return err
}

// GetAllUsers from repo.
func (r *Repo) GetAllUsers() ([]*model.User, error) {
	return nil, nil
}

// GetUser data from repo.
func (r *Repo) GetUser(id interface{}) (*model.User, error) {
	return &model.User{}, nil
}

// UpdatUser data in repo.
func (r *Repo) UpdateUser(*model.User) (*model.User, error) {
	return &model.User{}, nil
}

// DeletiUser data from repo.
func (r *Repo) DeleteUser(id interface{}) error {
	return nil
}

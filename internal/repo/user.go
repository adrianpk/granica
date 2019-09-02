package repo

import (
	"gitlab.com/mikrowezel/granica/pkg/auth"
)

// GetAllUsers from repo.
func (r *RepoHandler) GetAllUsers() ([]*auth.User, error) {
	return nil, nil
}

// CreateUser in repo.
func (r *RepoHandler) CreateUser(user *auth.User) (*auth.User, error) {
	tx, err := r.GetTx()
	if err != nil {
		return nil, err
	}

	st := `INSERT INTO users (id, slug, username, password_digest, email, given_name, middle_names, family_name, geolocation, :locale, base_tz, current_tz, starts_at, ends_at, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :username, :password_digest, :email, :given_name, :middle_names :family_name, :geolocation, :locale, :base_tz, :current_tz, :starts_at, :ends_at, :is_active, :is_deleted, created_by_id, updated_by_id, created_at, :updated_at)`

	_, err = tx.NamedExec(st, user)

	return user, err
}

// GetUser data from repo.
func (r *RepoHandler) GetUser(id interface{}) (*auth.User, error) {
	return &auth.User{}, nil
}

// UpdatUser data in repo.
func (r *RepoHandler) UpdateUser(*auth.User) (*auth.User, error) {
	return &auth.User{}, nil
}

// DeletiUser data from repo.
func (r *RepoHandler) DeleteUser(id interface{}) error {
	return nil
}

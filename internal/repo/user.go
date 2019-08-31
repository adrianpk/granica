package repo

import "gitlab.com/mikrowezel/granica/pkg/auth"

// GetAll users from repo.
func (r *RepoHandler) GetAll() ([]*auth.User, error) {
	return nil, nil
}

// Create user in repo.
func (r *RepoHandler) Create(user *auth.User) (*auth.User, error) {
	st := `INSERT INTO users (user_type, username, password_digest, email, first_name, last_name, created_at, updated_at)
	VALUES (:user_type, :username, :password_digest, :email, :first_name, :last_name, :created_at, :updated_at)`

	_, err := r.tx.NamedExec(st, user)

	return user, err
}

// Get user data from repo.
func (r *RepoHandler) GetUser(id interface{}) (*auth.User, error) {
	return &auth.User{}, nil
}

// Update user data in repo.
func (r *RepoHandler) UpdateUser(*auth.User) (*auth.User, error) {
	return &auth.User{}, nil
}

// Delete user data from repo.
func (r *RepoHandler) DeleteUser(id interface{}) error {
	return nil
}

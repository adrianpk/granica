package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"gitlab.com/mikrowezel/backend/config"
	logger "gitlab.com/mikrowezel/backend/log"
	"gitlab.com/mikrowezel/backend/granica/internal/model"
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

	st := `INSERT INTO users (id, tenant_id, slug, username, password_digest, email, given_name, middle_names, family_name, geolocation, locale, base_tz, current_tz, starts_at, ends_at, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :tenant_id, :slug, :username, :password_digest, :email, :given_name, :middle_names, :family_name, :geolocation, :locale, :base_tz, :current_tz, :starts_at, :ends_at, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	_, err := ur.Tx.NamedExec(st, user)

	return err
}

// GetAll users from repo.
func (ur *UserRepo) GetAll() (users []model.User, err error) {
	st := `SELECT * FROM users;`

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
func (ur *UserRepo) Update(user *model.User) error {
	ref, err := ur.Get(user.ID.String())
	if err != nil {
		return fmt.Errorf("cannot retrieve reference user: %s", err.Error())
	}

	user.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE users SET ")

	if user.Username.String != ref.Username.String {
		st.WriteString(strUpd("username", "username"))
		pcu = true
	}

	if user.PasswordDigest.String != ref.PasswordDigest.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("password_digest", "password_digest"))
		pcu = true
	}

	if user.Email.String != ref.Email.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("email", "email"))
		pcu = true
	}

	if user.GivenName.String != ref.GivenName.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("given_name", "given_name"))
		pcu = true
	}

	if user.MiddleNames.String != ref.MiddleNames.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("middle_names", "middle_names"))
		pcu = true
	}

	if user.FamilyName.String != ref.FamilyName.String {
		st.WriteString(preDelimiter(pcu))
		st.WriteString(strUpd("family_name", "family_name"))
		pcu = true
	}

	st.WriteString(" ")
	st.WriteString(whereID(ref.ID.String()))
	st.WriteString(";")

	//fmt.Println(st.String())

	_, err = ur.Tx.NamedExec(st.String(), user)

	return err
}

// Delete user from repo by ID.
func (ur *UserRepo) Delete(id string) error {
	st := `DELETE FROM USERS WHERE id = '%s';`
	st = fmt.Sprintf(st, id)

	_, err := ur.Tx.Exec(st)

	return err
}

// DeleteBySlug:w user from repo by slug.
func (ur *UserRepo) DeleteBySlug(slug string) error {
	st := `DELETE FROM USERS WHERE slug = '%s';`
	st = fmt.Sprintf(st, slug)

	_, err := ur.Tx.Exec(st)

	return err
}

// DeleteByusername user from repo by username.
func (ur *UserRepo) DeleteByUsername(username string) error {
	st := `DELETE FROM USERS WHERE username = '%s';`
	st = fmt.Sprintf(st, username)

	_, err := ur.Tx.Exec(st)

	return err
}

// preDelimiter selects a comma or space
// for each field in update statements.
func preDelimiter(upc bool) string {
	if upc {
		return ", "
	}
	return " "
}

// strUpdCol build an update colum fragment of type string.
func strUpd(colName, fieldName string) string {
	return fmt.Sprintf("%s = :%s", colName, fieldName)
}

// whereID build an SQL where clause for ID.
func whereID(id string) string {
	return fmt.Sprintf("WHERE id = '%s'", id)
}

// Commit transaction
func (ur *UserRepo) Commit() error {
	return ur.Tx.Commit()
}

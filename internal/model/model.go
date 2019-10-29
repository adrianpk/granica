package model

import (
	"database/sql"

	"github.com/lib/pq"
	"gitlab.com/mikrowezel/backend/db"
	m "gitlab.com/mikrowezel/backend/model"
	"golang.org/x/crypto/bcrypt"
)

// User -----------------------------------------------------------------------
type (
	// User model
	User struct {
		m.Identification
		Username          sql.NullString `db:"username" json:"username"`
		Password          string         `db:"password" json:"password"`
		PasswordDigest    sql.NullString `db:"password_digest" json:"-"`
		Email             sql.NullString `db:"email" json:"email"`
		EmailConfirmation sql.NullString `db:"emailConfirmation" json:"emailConfirmation"`
		GivenName         sql.NullString `db:"given_name" json:"givenName"`
		MiddleNames       sql.NullString `db:"middle_names" json:"middleNames"`
		FamilyName        sql.NullString `db:"family_name" json:"familyName"`
		Geolocation       db.NullPoint   `db:"geolocation" json:"geolocation"`
		Locale            sql.NullString `db:"locale" json:"locale"`
		BaseTZ            sql.NullString `db:"base_tz" json:"baseTZ"`
		CurrentTZ         sql.NullString `db:"current_tz" json:"currentTZ"`
		StartsAt          pq.NullTime    `db:"starts_at" json:"startsAt"`
		EndsAt            pq.NullTime    `db:"ends_at" json:"endsAt"`
		IsActive          sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted         sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		m.Audit
	}
)

// UpdatePasswordDigest if password changed.
func (user *User) UpdatePasswordDigest() (digest string, err error) {
	if user.Password == "" {
		return user.PasswordDigest.String, nil
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.PasswordDigest.String, err
	}
	user.PasswordDigest.String = string(hpass)
	return user.PasswordDigest.String, nil
}

// SetCreateValues sets de ID and slug.
func (user *User) SetCreateValues() error {
	pfx := user.Username.String
	user.Identification.SetCreateValues(pfx)
	user.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues
func (user *User) SetUpdateValues() error {
	user.Audit.SetUpdateValues()
	return nil
}

// Account --------------------------------------------------------------------
type (
	// Account model
	Account struct {
		m.Identification
		Name        sql.NullString `db:"name" json:"name"`
		OwnerID     sql.NullString `db:"owner_id" json:"ownerID"`
		ParentID    sql.NullString `db:"parent_id" json:"parentID"`
		AccountType sql.NullString `db:"account_type" json:"accountType"`
		Email       sql.NullString `db:"email" json:"email"`
		ShownName   sql.NullString `db:"shown_name" json:"shownName"`
		Geolocation db.NullPoint   `db:"geolocation" json:"geolocation"`
		Locale      sql.NullString `db:"locale" json:"locale"`
		BaseTZ      sql.NullString `db:"base_tz" json:"baseTZ"`
		CurrentTZ   sql.NullString `db:"current_tz" json:"currentTZ"`
		StartsAt    pq.NullTime    `db:"starts_at" json:"startsAt"`
		EndsAt      pq.NullTime    `db:"ends_at" json:"endsAt"`
		IsActive    sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted   sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		m.Audit
	}
)

// SetCreateValues sets de ID and slug.
func (account *Account) SetCreateValues() error {
	pfx := account.Name.String
	account.Identification.SetCreateValues(pfx)
	account.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues
func (account *Account) SetUpdateValues() error {
	account.Audit.SetUpdateValues()
	return nil
}

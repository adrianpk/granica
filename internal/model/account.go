package model

import (
	"database/sql"

	"github.com/lib/pq"
	m "gitlab.com/mikrowezel/backend/model"
)

type (
	// Account model
	Account struct {
		m.Identification
		OwnerID     sql.NullString `db:"owner_id" json:"ownerID"`
		ParentID    sql.NullString `db:"parent_id" json:"parentID"`
		AccountType sql.NullString `db:"account_type" json:"accountType"`
		Name        sql.NullString `db:"name" json:"name"`
		Email       sql.NullString `db:"email" json:"email"`
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

// Match condition for model.
func (account *Account) Match(tc *Account) bool {
	r := account.Identification.Match(tc.Identification) &&
		account.OwnerID == tc.OwnerID &&
		account.ParentID == tc.ParentID &&
		account.AccountType == tc.AccountType &&
		account.Name == tc.Name &&
		account.Email == tc.Email
	return r
}

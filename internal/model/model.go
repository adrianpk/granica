package model

import (
	"database/sql"

	"encoding/json"
	"time"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/mikrowezel/db"
	"golang.org/x/crypto/bcrypt"
)

type (
	// User model
	User struct {
		ID uuid.UUID `db:"id" json:"id"`
		Slug              sql.NullString `db:"slug" json:"slug"`
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
		CreatedByID       sql.NullString `db:"created_by_id" json:"createdByID"`
		UpdatedByID       sql.NullString `db:"updated_by_id" json:"updatedByID"`
		CreatedAt         pq.NullTime    `db:"created_at" json:"createdAt"`
		UpdatedAt         pq.NullTime    `db:"updated_at" json:"updatedAt"`
	}
)

// GetID representation.
func (user *User) GetID() interface{} {
	return user.ID.Value
}

// SetID for user.
func (user *User) SetID(id uuid.UUID) {
	user.ID = id
}

// GenID for user.
func (user *User) GenID() {
	if user.ID == uuid.Nil {
		user.ID = uuid.NewV4()
	}
}

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

// MarshalJSON - Custom MarshalJSON function.
func (user *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		CreatedAt int64 `json:"createdAt"`
		UpdatedAt int64 `json:"updatedAt"`
	}{
		Alias:     (*Alias)(user),
		CreatedAt: user.CreatedAt.Time.Unix(),
		UpdatedAt: user.UpdatedAt.Time.Unix(),
	})
}

// UnmarshalJSON - Custom UnmarshalJSON function.
func (user *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		*Alias
		CreatedAt int64 `json:"createdAt"`
		UpdatedAt int64 `json:"updatedAt"`
	}{
		Alias:     (*Alias)(user),
		CreatedAt: user.CreatedAt.Time.Unix(),
		UpdatedAt: user.UpdatedAt.Time.Unix(),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	tc := time.Unix(aux.CreatedAt, 0)
	tu := time.Unix(aux.UpdatedAt, 0)
	user.CreatedAt = pq.NullTime{tc, true}
	user.UpdatedAt = pq.NullTime{tu, true}
	return nil
}

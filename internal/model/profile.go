package model

import (
	"database/sql"

	"github.com/lib/pq"
	"gitlab.com/mikrowezel/backend/db"
	m "gitlab.com/mikrowezel/backend/model"
)

type (
	// Profile model
	Profile struct {
		m.Identification
		OwnerID        sql.NullString `db:"owner_id" json:"ownerID"`
		ProfileType    sql.NullString `db:"profile_type" json:"profileType"`
		Name           sql.NullString `db:"name" json:"name"`
		Email          sql.NullString `db:"email" json:"email"`
		Description    sql.NullString `db:"description" json:"description"`
		Location       sql.NullString `db:"location" json:"location"`
		Bio            sql.NullString `db:"bio" json:"bio"`
		Moto           sql.NullString `db:"moto" json:"moto"`
		Website        sql.NullString `db:"website" json:"website"`
		AniversaryDate pq.NullTime    `db:"aniversary_date" json:"aniversaryDate"`
		AvatarPath     sql.NullString `db:"avatar_path" json:"avatarPath"`
		HeaderPath     sql.NullString `db:"header_path" json:"headerPath"`
		Geolocation    db.NullPoint   `db:"geolocation" json:"geolocation"`
		Locale         sql.NullString `db:"locale" json:"locale"`
		BaseTZ         sql.NullString `db:"base_tz" json:"baseTZ"`
		CurrentTZ      sql.NullString `db:"current_tz" json:"currentTZ"`
		IsActive       sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted      sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		m.Audit
	}
)

// SetCreateValues sets de ID and slug.
func (profile *Profile) SetCreateValues() error {
	pfx := profile.Name.String
	profile.Identification.SetCreateValues(pfx)
	profile.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues
func (profile *Profile) SetUpdateValues() error {
	profile.Audit.SetUpdateValues()
	return nil
}

// Match condition for model.
func (profile *Profile) Match(tc *Profile) bool {
	r := profile.Identification.Match(tc.Identification) &&
		profile.OwnerID == tc.OwnerID &&
		profile.ProfileType == tc.ProfileType &&
		profile.Name == tc.Name &&
		profile.Email == tc.Email &&
		profile.Description == tc.Description &&
		profile.Location == tc.Location &&
		profile.Bio == tc.Bio &&
		profile.Moto == tc.Moto &&
		profile.Website == tc.Website &&
		profile.AniversaryDate == tc.AniversaryDate &&
		profile.AvatarPath == tc.AvatarPath &&
		profile.HeaderPath == tc.HeaderPath
	return r
}

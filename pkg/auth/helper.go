package auth

import (
	"gitlab.com/mikrowezel/db"
	"gitlab.com/mikrowezel/granica/internal/model"
)

// toModel creates a User model fron transport values.
func (cur *CreateUserReq) toModel() model.User {
	return model.User{
		Username:          db.ToNullString(cur.Username),
		Password:          cur.Password,
		Email:             db.ToNullString(cur.Email),
		EmailConfirmation: db.ToNullString(cur.EmailConfirmation),
		GivenName:         db.ToNullString(cur.GivenName),
		MiddleNames:       db.ToNullString(cur.MiddleNames),
		FamilyName:        db.ToNullString(cur.FamilyName),
		// Geolocation:    db.ToNullGeometry(cur.Lat, cur.Lng)
	}
}

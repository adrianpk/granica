package auth

import "gitlab.com/mikrowezel/granica/internal/model"

// toModel creates a User model fron transport values.
func (cur *CreateUserReq) toModel() model.User {
	return model.User{
		Username:          toNullString(cur.Username),
		Password:          cur.Password,
		Email:             toNullString(cur.Email),
		EmailConfirmation: toNullString(cur.EmailConfirmation),
		GivenName:         toNullString(cur.GivenName),
		MiddleNames:       toNullString(cur.MiddleNames),
		FamilyName:        toNullString(cur.FamilyName),
		// Geolocation:    toNullGeometry(cur.Lat, cur.Lng)
	}
}

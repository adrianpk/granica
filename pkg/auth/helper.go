package auth

import (
	"encoding/json"
	"fmt"

	"gitlab.com/mikrowezel/backend/db"
	"gitlab.com/mikrowezel/granica/internal/model"
)

// createUser ------------------------------------------------------------------
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

func (cur *CreateUserRes) fromModel(u *model.User, msg string, err error) {
	if u != nil {
		cur.User = User{
			Slug:        u.Slug.String,
			Username:    u.Username.String,
			Password:    "",
			Email:       u.Email.String,
			GivenName:   u.GivenName.String,
			MiddleNames: u.MiddleNames.String,
			FamilyName:  u.FamilyName.String,
			Lat:         fmt.Sprintf("%f", u.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", u.Geolocation.Point.Lng),
		}
	}
	cur.Msg = msg
	if err != nil {
		cur.Error = err.Error()
	}
}

// getUsers -------------------------------------------------------------------
func (gur *GetUsersRes) fromModel(us []model.User, msg string, err error) {
	gurUsers := []User{}
	for _, u := range us {
		gur := User{
			Username:    u.Username.String,
			Password:    "",
			Email:       u.Email.String,
			GivenName:   u.GivenName.String,
			MiddleNames: u.MiddleNames.String,
			FamilyName:  u.FamilyName.String,
			Lat:         fmt.Sprintf("%f", u.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", u.Geolocation.Point.Lng),
		}
		gurUsers = append(gurUsers, gur)
	}
	gur.Users = gurUsers
	gur.Msg = msg
	if err != nil {
		gur.Error = err.Error()
	}
}

// getUser ---------------------------------------------------------------------
func (gur *GetUserReq) toModel() model.User {
	return model.User{
		Username: db.ToNullString(gur.Identifier.Username),
	}
}

func (cur *GetUserRes) fromModel(u *model.User, msg string, err error) {
	if u != nil {
		cur.User = User{
			Username:    u.Username.String,
			Password:    "",
			Email:       u.Email.String,
			GivenName:   u.GivenName.String,
			MiddleNames: u.MiddleNames.String,
			FamilyName:  u.FamilyName.String,
			Lat:         fmt.Sprintf("%f", u.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", u.Geolocation.Point.Lng),
		}
	}
	cur.Msg = msg
	if err != nil {
		cur.Error = err.Error()
	}
}

// updateUser ------------------------------------------------------------------
func (a *Auth) makeUpdateUserResJSON(u *model.User, msg string, err error) ([]byte, error) {
	cur := UpdateUserRes{}
	cur.fromModel(u, msg, err)
	return a.toJSON(cur)
}

// toModel creates a User model from transport values.
func (cur *UpdateUserReq) toModel() model.User {
	return model.User{
		Username:          db.ToNullString(cur.User.Username),
		Password:          cur.Password,
		Email:             db.ToNullString(cur.Email),
		EmailConfirmation: db.ToNullString(cur.EmailConfirmation),
		GivenName:         db.ToNullString(cur.GivenName),
		MiddleNames:       db.ToNullString(cur.MiddleNames),
		FamilyName:        db.ToNullString(cur.FamilyName),
		// Geolocation:    db.ToNullGeometry(cur.Lat, cur.Lng)
	}
}

func (uur *UpdateUserRes) fromModel(u *model.User, msg string, err error) {
	if u != nil {
		uur.User = User{
			Slug:        u.Slug.String,
			Username:    u.Username.String,
			Password:    "",
			Email:       u.Email.String,
			GivenName:   u.GivenName.String,
			MiddleNames: u.MiddleNames.String,
			FamilyName:  u.FamilyName.String,
			Lat:         fmt.Sprintf("%f", u.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", u.Geolocation.Point.Lng),
		}
	}
	uur.Msg = msg
	if err != nil {
		uur.Error = err.Error()
	}
}

// deleteUser ------------------------------------------------------------------
func (dur *DeleteUserRes) fromModel(u *model.User, msg string, err error) {
	dur.Msg = msg
	if err != nil {
		dur.Error = err.Error()
	}
}

// TODO: Move to response method.
func (a *Auth) toJSON(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

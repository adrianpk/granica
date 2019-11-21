package transport

import (
	"fmt"

	"gitlab.com/mikrowezel/backend/db"

	"gitlab.com/mikrowezel/backend/granica/internal/model"
)

func (req *CreateUserReq) ToModel() model.User {
	return model.User{
		Username:          db.ToNullString(req.Username),
		Password:          req.Password,
		Email:             db.ToNullString(req.Email),
		EmailConfirmation: db.ToNullString(req.EmailConfirmation),
		GivenName:         db.ToNullString(req.GivenName),
		MiddleNames:       db.ToNullString(req.MiddleNames),
		FamilyName:        db.ToNullString(req.FamilyName),
		// Geolocation:    db.ToNullGeometry(req.Lat, req.Lng)
	}
}

func (res *CreateUserRes) FromModel(m *model.User, msg string, err error) {
	if m != nil {
		res.User = User{
			Slug:        m.Slug.String,
			Username:    m.Username.String,
			Password:    "",
			Email:       m.Email.String,
			GivenName:   m.GivenName.String,
			MiddleNames: m.MiddleNames.String,
			FamilyName:  m.FamilyName.String,
			Lat:         fmt.Sprintf("%f", m.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", m.Geolocation.Point.Lng),
			IsNew:       m.IsNew(),
		}
	}
}

func (res *GetUsersRes) FromModel(ms []model.User, msg string, err error) {
	resUsers := []User{}
	for _, m := range ms {
		res := User{
			Username:    m.Username.String,
			Password:    "",
			Email:       m.Email.String,
			GivenName:   m.GivenName.String,
			MiddleNames: m.MiddleNames.String,
			FamilyName:  m.FamilyName.String,
			Lat:         fmt.Sprintf("%f", m.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", m.Geolocation.Point.Lng),
		}
		resUsers = append(resUsers, res)
	}
	res.Users = resUsers
}

func (req *GetUserReq) ToModel() model.User {
	return model.User{
		Username: db.ToNullString(req.Identifier.Username),
	}
}

func (res *GetUserRes) FromModel(m *model.User, msg string, err error) {
	if m != nil {
		res.User = User{
			Username:    m.Username.String,
			Password:    "",
			Email:       m.Email.String,
			GivenName:   m.GivenName.String,
			MiddleNames: m.MiddleNames.String,
			FamilyName:  m.FamilyName.String,
			Lat:         fmt.Sprintf("%f", m.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", m.Geolocation.Point.Lng),
		}
	}
}

//func (a *Auth) makeUpdateUserResJSON(m *model.User, msg string, err error) ([]byte, error) {
//res := UpdateUserRes{}
//res.FromModel(m, msg, err)
//return a.toJSON(res.User)
//}

// ToModel creates a User model from transport values.
func (req *UpdateUserReq) ToModel() model.User {
	return model.User{
		Username:          db.ToNullString(req.User.Username),
		Password:          req.Password,
		Email:             db.ToNullString(req.Email),
		EmailConfirmation: db.ToNullString(req.EmailConfirmation),
		GivenName:         db.ToNullString(req.GivenName),
		MiddleNames:       db.ToNullString(req.MiddleNames),
		FamilyName:        db.ToNullString(req.FamilyName),
		// Geolocation:    db.ToNullGeometry(req.Lat, req.Lng)
	}
}

func (res *UpdateUserRes) FromModel(m *model.User, msg string, err error) {
	if m != nil {
		res.User = User{
			Slug:        m.Slug.String,
			Username:    m.Username.String,
			Password:    "",
			Email:       m.Email.String,
			GivenName:   m.GivenName.String,
			MiddleNames: m.MiddleNames.String,
			FamilyName:  m.FamilyName.String,
			Lat:         fmt.Sprintf("%f", m.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", m.Geolocation.Point.Lng),
		}
	}
}

func (res *DeleteUserRes) FromModel(m *model.User, msg string, err error) {
}

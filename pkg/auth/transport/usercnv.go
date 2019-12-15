package transport

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/mikrowezel/backend/db"
	m "gitlab.com/mikrowezel/backend/model"
	"gitlab.com/mikrowezel/backend/web"

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

func (res *CreateUserRes) FromForm(cs *web.Form) {
	res.User = User{
		Username:    cs.Get("username"),
		Password:    cs.Get("password"),
		Email:       cs.Get("email"),
		GivenName:   cs.Get("given-name"),
		MiddleNames: cs.Get("middle-names"),
		FamilyName:  cs.Get("family-name"),
	}

	res.Errors = cs.Errors
}

func (res *IndexUsersRes) FromModel(ms []model.User, msg string, err error) {
	resUsers := []User{}
	for _, m := range ms {
		res := User{
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
		resUsers = append(resUsers, res)
	}
	res.Users = resUsers
}

func (req *GetUserReq) ToModel() model.User {
	return model.User{
		Identification: m.Identification{
			ID:       uuid.UUID{},
			TenantID: db.ToNullString(""),
			Slug:     db.ToNullString(req.Identifier.Slug),
		},
	}
}

func (res *GetUserRes) FromModel(m *model.User, msg string, err error) {
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

// ToModel creates a User model from transport values.
func (req *UpdateUserReq) ToModel() model.User {
	return model.User{
		Identification: m.Identification{
			ID:       uuid.UUID{},
			TenantID: db.ToNullString(""),
			Slug:     db.ToNullString(req.User.Slug),
		},
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

func (res *UpdateUserRes) FromForm(cs *web.Form) {
	res.User = User{
		Username:    cs.Get("username"),
		Password:    cs.Get("password"),
		Email:       cs.Get("email"),
		GivenName:   cs.Get("given-name"),
		MiddleNames: cs.Get("middle-names"),
		FamilyName:  cs.Get("family-name"),
	}

	res.Errors = cs.Errors
}

func (res *DeleteUserRes) FromModel(m *model.User, msg string, err error) {
}

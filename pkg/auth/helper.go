package auth

import (
	"fmt"

	"gitlab.com/mikrowezel/backend/db"
	"gitlab.com/mikrowezel/granica/internal/model"
)

// User -----------------------------------------------------------------------
func (req *CreateUserReq) toModel() model.User {
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

func (res *CreateUserRes) fromModel(m *model.User, msg string, err error) {
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
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

func (res *GetUsersRes) fromModel(ms []model.User, msg string, err error) {
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
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

func (req *GetUserReq) toModel() model.User {
	return model.User{
		Username: db.ToNullString(req.Identifier.Username),
	}
}

func (res *GetUserRes) fromModel(m *model.User, msg string, err error) {
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
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

func (a *Auth) makeUpdateUserResJSON(m *model.User, msg string, err error) ([]byte, error) {
	res := UpdateUserRes{}
	res.fromModel(m, msg, err)
	return a.toJSON(res.User)
}

// toModel creates a User model from transport values.
func (req *UpdateUserReq) toModel() model.User {
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

func (res *UpdateUserRes) fromModel(m *model.User, msg string, err error) {
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
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

func (res *DeleteUserRes) fromModel(m *model.User, msg string, err error) {
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

// Account -----------------------------------------------------------------------
func (req *CreateAccountReq) toModel() model.Account {
	return model.Account{
		Name:        db.ToNullString(req.Name),
		AccountType: db.ToNullString(req.AccountType),
		OwnerID:     db.ToNullString(req.OwnerID),
		ParentID:    db.ToNullString(req.ParentID),
		Email:       db.ToNullString(req.Email),
		ShownName:   db.ToNullString(req.ShownName),
		// Geolocation:    db.ToNullGeometry(req.Lat, req.Lng)
	}
}

func (res *CreateAccountRes) fromModel(m *model.Account, msg string, err error) {
	if m != nil {
		res.Account = Account{
			Slug:        m.Slug.String,
			Name:        m.Name.String,
			AccountType: m.AccountType.String,
			OwnerID:     m.OwnerID.String,
			ParentID:    m.ParentID.String,
			Email:       m.Email.String,
			ShownName:   m.ShownName.String,
			Lat:         fmt.Sprintf("%f", m.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", m.Geolocation.Point.Lng),
		}
	}
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

// getAccounts -------------------------------------------------------------------
func (res *GetAccountsRes) fromModel(ms []model.Account, msg string, err error) {
	resAccounts := []Account{}
	for _, m := range ms {
		res := Account{
			Name:        m.Name.String,
			AccountType: m.AccountType.String,
			OwnerID:     m.OwnerID.String,
			ParentID:    m.ParentID.String,
			Email:       m.Email.String,
			ShownName:   m.ShownName.String,
			Lat:         fmt.Sprintf("%f", m.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", m.Geolocation.Point.Lng),
		}
		resAccounts = append(resAccounts, res)
	}
	res.Accounts = resAccounts
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

// getAccount ---------------------------------------------------------------------
func (req *GetAccountReq) toModel() model.Account {
	return model.Account{
		Identification: model.Identification{
			Slug: db.ToNullString(req.Identifier.Slug),
		},
	}
}

func (res *GetAccountRes) fromModel(m *model.Account, msg string, err error) {
	if m != nil {
		res.Account = Account{
			Name:        m.Name.String,
			AccountType: m.AccountType.String,
			OwnerID:     m.OwnerID.String,
			ParentID:    m.ParentID.String,
			Email:       m.Email.String,
			ShownName:   m.ShownName.String,
			Lat:         fmt.Sprintf("%f", m.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", m.Geolocation.Point.Lng),
		}
	}
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

// updateAccount ------------------------------------------------------------------
func (a *Auth) makeUpdateAccountResJSON(m *model.Account, msg string, err error) ([]byte, error) {
	res := UpdateAccountRes{}
	res.fromModel(m, msg, err)
	return a.toJSON(res.Account)
}

// toModel creates a Account model from transport values.
func (req *UpdateAccountReq) toModel() model.Account {
	return model.Account{
		Name:        db.ToNullString(req.Name),
		AccountType: db.ToNullString(req.AccountType),
		OwnerID:     db.ToNullString(req.OwnerID),
		ParentID:    db.ToNullString(req.ParentID),
		Email:       db.ToNullString(req.Email),
		ShownName:   db.ToNullString(req.ShownName),
		// Geolocation:    db.ToNullGeometry(req.Lat, req.Lng)
	}
}

func (res *UpdateAccountRes) fromModel(m *model.Account, msg string, err error) {
	if m != nil {
		res.Account = Account{
			Slug:        m.Slug.String,
			Name:        m.Name.String,
			AccountType: m.AccountType.String,
			OwnerID:     m.OwnerID.String,
			ParentID:    m.ParentID.String,
			Email:       m.Email.String,
			ShownName:   m.ShownName.String,
			Lat:         fmt.Sprintf("%f", m.Geolocation.Point.Lat),
			Lng:         fmt.Sprintf("%f", m.Geolocation.Point.Lng),
		}
	}
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

// deleteAccount ------------------------------------------------------------------
func (res *DeleteAccountRes) fromModel(m *model.Account, msg string, err error) {
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

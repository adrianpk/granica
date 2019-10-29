package transport

import (
	"fmt"

	"gitlab.com/mikrowezel/backend/db"
	"gitlab.com/mikrowezel/backend/granica/internal/model"
	m "gitlab.com/mikrowezel/backend/model"
)

func (req *CreateAccountReq) ToModel() model.Account {
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

func (res *CreateAccountRes) FromModel(m *model.Account, msg string, err error) {
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
func (res *GetAccountsRes) FromModel(ms []model.Account, msg string, err error) {
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
func (req *GetAccountReq) ToModel() model.Account {
	return model.Account{
		Identification: m.Identification{
			Slug: db.ToNullString(req.Identifier.Slug),
		},
	}
}

func (res *GetAccountRes) FromModel(m *model.Account, msg string, err error) {
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
//func (a *Auth) makeUpdateAccountResJSON(m *model.Account, msg string, err error) ([]byte, error) {
//res := UpdateAccountRes{}
//res.FromModel(m, msg, err)
//return a.toJSON(res.Account)
//}

// ToModel creates a Account model from transport values.
func (req *UpdateAccountReq) ToModel() model.Account {
	return model.Account{
		Identification: m.Identification{
			Slug: db.ToNullString(req.Identifier.Slug),
		},
		Name:        db.ToNullString(req.Name),
		AccountType: db.ToNullString(req.AccountType),
		OwnerID:     db.ToNullString(req.OwnerID),
		ParentID:    db.ToNullString(req.ParentID),
		Email:       db.ToNullString(req.Email),
		ShownName:   db.ToNullString(req.ShownName),
		// Geolocation:    db.ToNullGeometry(req.Lat, req.Lng)
	}
}

func (res *UpdateAccountRes) FromModel(m *model.Account, msg string, err error) {
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
func (res *DeleteAccountRes) FromModel(m *model.Account, msg string, err error) {
	res.Msg = msg
	if err != nil {
		res.Error = err.Error()
	}
}

package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/mikrowezel/backend/db"
	"gitlab.com/mikrowezel/granica/internal/model"
)

// createUserResponse creates a CreateUserRes, encodes it to JSON
// and write it using the ResponseWriter.
func (a *Auth) createUserResponse(w http.ResponseWriter, r *http.Request, u *model.User, msg string, err error) error {
	out, err := a.makeCreateUserResJSON(u, msg, err)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
	return nil
}

// makeCreateUserResJSON creates a JSON output using user model and error.
func (a *Auth) makeCreateUserResJSON(u *model.User, msg string, err error) ([]byte, error) {
	cur := CreateUserRes{}
	cur.fromModel(u, msg, err)
	return cur.toJSON()
}

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

// fromModel update CreateUserRes using model values.
func (cur *CreateUserRes) fromModel(u *model.User, msg string, err error) {
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

func (cur *CreateUserRes) toJSON() ([]byte, error) {
	return json.Marshal(cur)
}

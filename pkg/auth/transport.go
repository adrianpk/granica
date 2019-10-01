package auth

import "encoding/json"

type (
	// User request and response data.
	User struct {
		Username          string `json:"username"`
		Password          string `json:"password"`
		Email             string `json:"email"`
		EmailConfirmation string `json:"emailConfirmation"`
		GivenName         string `json:"givenName"`
		MiddleNames       string `json:"middleNames"`
		FamilyName        string `json:"familyName"`
		Lat               string `json:"geolocation"`
		Lng               string `json:"geolocation"`
	}

	// CreateUserReq input data.
	CreateUserReq struct {
		User
	}

	// CreateUserRes output data.
	CreateUserRes struct {
		User
		Msg   string `json: msg,omitempty`
		Error string `json:"err,omitempty"`
	}
)

func toJSON(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

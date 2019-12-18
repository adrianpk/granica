package transport

import (
	"gitlab.com/mikrowezel/backend/service"
	"gitlab.com/mikrowezel/backend/web"
)

type (
	// User request and response data.
	User struct {
		Slug              string `json:"slug" schema:"slug"`
		Username          string `json:"username" schema:"username"`
		Password          string `json:"password" schema:"password"`
		Email             string `json:"email" schema:"email"`
		EmailConfirmation string `json:"emailConfirmation" schema:"email-confirmation"`
		GivenName         string `json:"givenName" schema:"given-name"`
		MiddleNames       string `json:"middleNames" schema:"middle-names"`
		FamilyName        string `json:"familyName" schema:"family-name"`
		LastIP            string `json:"lastIP" schema:"last-ip"`
		VerifyToken       string `json:"verifyToken" schema:"verify-token"`
		IsVerified        bool   `json:"isVerified" schema:"is-verified"`
		Lat               string `json:"lat" schema: "lat"`
		Lng               string `json:"lng" schema: "lng"`
		IsNew             bool
	}

	// SignIn
	SignIn struct {
		Username string `json:"username" schema:"username"`
		Password string `json:"password" schema:"password"`
	}

	Users []User
)

func (u User) GetSlug() string {
	return u.Slug
}

type (
	// CreateUserReq input data.
	CreateUserReq struct {
		User
	}

	// CreateUserRes output data.
	CreateUserRes struct {
		User
		Action web.Action
		Errors service.ErrorSet
	}
)

type (
	// IndexUsersReq input data.
	IndexUsersReq struct {
	}

	// IndexUsersRes output data.
	IndexUsersRes struct {
		Users
	}
)

type (
	// GetUserReq input data.
	GetUserReq struct {
		Identifier
	}

	// GetUserRes output data.
	GetUserRes struct {
		User
		Action web.Action
		Errors service.ErrorSet
	}
)

type (
	// UpdateUserReq input data.
	UpdateUserReq struct {
		Identifier
		User
	}

	// UpdateUserRes output data.
	UpdateUserRes struct {
		User
		Action web.Action
		Errors service.ErrorSet
	}
)

type (
	// DeleteUserReq input data.
	DeleteUserReq struct {
		Identifier
	}

	// DeleteUserRes output data.
	DeleteUserRes struct {
	}
)

type (
	// SignUpUserReq input data.
	SignUpUserReq struct {
		User
	}

	// SignUpUserRes output data.
	SignUpUserRes struct {
		User
		Action web.Action
		Errors service.ErrorSet
	}
)

type (
	// SignInUserReq input data.
	SignInUserReq struct {
		SignIn
	}

	// SignInUserRes output data.
	SignInUserRes struct {
		User
		Action web.Action
		Errors service.ErrorSet
	}
)

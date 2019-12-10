package transport

import (
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
		Lat               string `json:"lat" schema: "lat"`
		Lng               string `json:"lng" schema: "lng"`
		IsNew             bool
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

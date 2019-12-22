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
		ConfirmationToken string `json:"confirmationToken" schema:"verify-token"`
		IsConfirmed       bool   `json:"isConfirmed" schema:"is-confirmed"`
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
		// Action can be used to reuse form templates letting change target and method from controller.
		Action web.Action
		// Errors stores localizable errors message IDs for model properties.
		// Mainly used to show messages on fields with errors after validation.
		Errors service.ErrorSet
		// Msg stores localizable message ID for the whole model.
		// Mainly used to show a message related to current state of the model.
		MsgID string
		// Mainly used for debugging porpouses and/or to show relevant info to admin users.
		err error
	}
)

type (
	// IndexUsersReq input data.
	IndexUsersReq struct {
	}

	// IndexUsersRes output data.
	IndexUsersRes struct {
		Users
		// Action can be used to reuse form templates letting change target and method from controller.
		Action web.Action
		// Errors stores localizable errors message IDs for model properties.
		// Mainly used to show messages on fields with errors after validation.
		Errors service.ErrorSet
		// Msg stores localizable message ID for the whole model.
		// Mainly used to show a message related to current state of the model.
		MsgID string
		// Mainly used for debugging porpouses and/or to show relevant info to admin users.
		err error
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
		// Action can be used to reuse form templates letting change target and method from controller.
		Action web.Action
		// Errors stores localizable errors message IDs for model properties.
		// Mainly used to show messages on fields with errors after validation.
		Errors service.ErrorSet
		// Msg stores localizable message ID for the whole model.
		// Mainly used to show a message related to current state of the model.
		MsgID string
		// Mainly used for debugging porpouses and/or to show relevant info to admin users.
		err error
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
		// Action can be used to reuse form templates letting change target and method from controller.
		Action web.Action
		// Errors stores localizable errors message IDs for model properties.
		// Mainly used to show messages on fields with errors after validation.
		Errors service.ErrorSet
		// MsgID stores localizable message ID for the whole model.
		// Mainly used to show a message related to current state of the model.
		MsgID string
		// Mainly used for debugging porpouses and/or to show relevant info to admin users.
		err error
	}
)

type (
	// DeleteUserReq input data.
	DeleteUserReq struct {
		Identifier
	}

	// DeleteUserRes output data.
	DeleteUserRes struct {
		// MsgID stores localizable message ID for the whole model.
		// Mainly used to show a message related to current state of the model.
		MsgID string
		// Mainly used for debugging porpouses and/or to show relevant info to admin users.
		err error
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
		// Action can be used to reuse form templates letting change target and method from controller.
		Action web.Action
		// Errors stores localizable errors message IDs for model properties.
		// Mainly used to show messages on fields with errors after validation.
		Errors service.ErrorSet
		// MsgID stores localizable message ID for the whole model.
		// Mainly used to show a message related to current state of the model.
		MsgID string
		// Mainly used for debugging porpouses and/or to show relevant info to admin users.
		err error
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
		// Action can be used to reuse form templates letting change target and method from controller.
		Action web.Action
		// Errors stores localizable errors message IDs for model properties.
		// Mainly used to show messages on fields with errors after validation.
		Errors service.ErrorSet
		// MsgID stores localizable message ID for the whole model.
		// Mainly used to show a message related to current state of the model.
		MsgID string
		// Mainly used for debugging porpouses and/or to show relevant info to admin users.
		err error
	}
)

package transport

type (
	// User request and response data.
	User struct {
		Slug              string `json:"slug"`
		Username          string `json:"username"`
		Password          string `json:"password"`
		Email             string `json:"email"`
		EmailConfirmation string `json:"emailConfirmation"`
		GivenName         string `json:"givenName"`
		MiddleNames       string `json:"middleNames"`
		FamilyName        string `json:"familyName"`
		Lat               string `json:"lat"`
		Lng               string `json:"lng"`
	}

	Users []User
)

func (u *User) GetSlug() string {
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
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

type (
	// GetUsersReq input data.
	GetUsersReq struct {
	}

	// GetUsersRes output data.
	GetUsersRes struct {
		Users
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
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
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
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
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

type (
	// DeleteUserReq input data.
	DeleteUserReq struct {
		Identifier
	}

	// DeleteUserRes output data.
	DeleteUserRes struct {
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

package auth

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
		Lat               string `json:"lat"`
		Lng               string `json:"lng"`
	}

	Users []User
)

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
	// GetUsersRes output data.
	GetUsersRes struct {
		Users
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

type (
	// GetUserRes output data.
	GetUserRes struct {
		User
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

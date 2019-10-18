package auth

// User -----------------------------------------------------------------------
type (
	Identifier struct {
		Slug     string
		Username string
	}

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

// Account --------------------------------------------------------------------
type (
	// Account request and response data.
	Account struct {
		TenantID    string `json:"tenantID"`
		Slug        string `json:"slug"`
		Name        string `json:"name"`
		OwnerID     string `json:"ownerID"`
		ParentID    string `json:"parentID"`
		AccountType string `json:"accountType"`
		Email       string `json:"email"`
		ShownName   string `json:"shownName"`
		Lat         string `json:"lat"`
		Lng         string `json:"lng"`
		StartsAt    string `json:"startsAt"`
		EndsAt      string `json:"endsAt"`
	}

	Accounts []Account
)

type (
	// CreateAccountReq input data.
	CreateAccountReq struct {
		Account
	}

	// CreateAccountRes output data.
	CreateAccountRes struct {
		Account
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

type (
	// GetAccountsReq input data.
	GetAccountsReq struct {
	}

	// GetAccountsRes output data.
	GetAccountsRes struct {
		Accounts
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

type (
	// GetAccountReq input data.
	GetAccountReq struct {
		Identifier
	}

	// GetAccountRes output data.
	GetAccountRes struct {
		Account
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

type (
	// UpdateAccountReq input data.
	UpdateAccountReq struct {
		Identifier
		Account
	}

	// UpdateAccountRes output data.
	UpdateAccountRes struct {
		Account
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

type (
	// DeleteAccountReq input data.
	DeleteAccountReq struct {
		Identifier
	}

	// DeleteAccountRes output data.
	DeleteAccountRes struct {
		Msg   string `json:"msg,omitempty"`
		Error string `json:"err,omitempty"`
	}
)

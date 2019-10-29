package transport

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

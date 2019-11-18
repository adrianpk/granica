package web

import (
	"net/http"

	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
	"gitlab.com/mikrowezel/backend/web"
)

const (
	userRes = "user"
)

const (
	UserCtxKey web.ContextKey = "user"
)

func (ep *Endpoint) InitCreateUser(w http.ResponseWriter, r *http.Request) {
	ep.Redirect(w, r, "http://google.com")
}

func (ep *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req tp.CreateUserReq
	var res tp.CreateUserRes

	// Service
	err := ep.Service().CreateUser(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.Redirect(w, r, "/")
		return
	}
}

// GetUsers web endpoint.
func (ep *Endpoint) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req tp.GetUsersReq
	var res tp.GetUsersRes

	// Service
	err := ep.Service().GetUsers(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.Redirect(w, r, "/")
		return
	}

	// Wrap response
	wr := ep.OKRes(res)

	// Template
	ts, err := ep.TemplateFor(userRes, web.IndexTmpl)
	if err != nil {
		ep.Redirect(w, r, "/")
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.Log().Error(err)
		ep.Redirect(w, r, "/")
	}
}

// GetUser web endpoint.
func (ep *Endpoint) GetUser(w http.ResponseWriter, r *http.Request) {
}

// UpdateUser web endpoint.
func (ep *Endpoint) UpdateUser(w http.ResponseWriter, r *http.Request) {
}

// DeleteUser web endpoint.
func (ep *Endpoint) DeleteUser(w http.ResponseWriter, r *http.Request) {
}

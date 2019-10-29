package jsonrest

import (
	"encoding/json"
	"errors"
	"net/http"

	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
)

const (
	UserCtxKey contextKey = "user"
)

func (ep *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req tp.CreateUserReq
	var res tp.CreateUserRes

	// Decode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Service
	err = ep.service.CreateUser(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

func (ep *Endpoint) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req tp.GetUsersReq
	var res tp.GetUsersRes

	// Service
	err := ep.service.GetUsers(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

func (ep *Endpoint) GetUser(w http.ResponseWriter, r *http.Request) {
	var req tp.GetUserReq
	var res tp.GetUserRes

	ctx := r.Context()
	username, ok := ctx.Value(UserCtxKey).(string)
	if !ok {
		e := errors.New("invalid username")
		ep.Log().Error(e)
		ep.writeResponse(w, res)
		return
	}

	// Service
	req.Username = username
	err := ep.service.GetUser(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

func (ep *Endpoint) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req tp.UpdateUserReq
	var res tp.UpdateUserRes

	ctx := r.Context()
	username, ok := ctx.Value(UserCtxKey).(string)
	if !ok {
		e := errors.New("invalid username")
		ep.Log().Error(e)
		ep.writeResponse(w, res)
		return
	}

	// Decode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Service
	req.Identifier.Username = username
	err = ep.service.UpdateUser(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

func (ep *Endpoint) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req tp.DeleteUserReq
	var res tp.DeleteUserRes

	ctx := r.Context()
	username, ok := ctx.Value(UserCtxKey).(string)
	if !ok {
		e := errors.New("invalid username")
		ep.Log().Error(e)
		ep.writeResponse(w, res)
		return
	}

	// Service
	req.Identifier.Username = username
	err := ep.service.DeleteUser(req, &res)
	if err != nil {
		e := errors.New("invalid username")
		ep.Log().Error(e)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

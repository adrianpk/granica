package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	createErr = "Cannot create entity"
	getAllErr = "Cannot get entity"
	getErr    = "Cannot get entity"
	updateErr = "Cannot update entity"
	deleteErr = "Cannot delete entity"
)

// User -----------------------------------------------------------------------
func (a *Auth) CreateUserJSON(w http.ResponseWriter, r *http.Request) {
	var req CreateUserReq
	var res CreateUserRes

	// Decode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Service
	err = a.CreateUser(req, &res)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

func (a *Auth) GetUsersJSON(w http.ResponseWriter, r *http.Request) {
	var req GetUsersReq
	var res GetUsersRes

	// Service
	err := a.GetUsers(req, &res)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

func (a *Auth) GetUserJSON(w http.ResponseWriter, r *http.Request) {
	var req GetUserReq
	var res GetUserRes

	ctx := r.Context()
	username, ok := ctx.Value(userCtxKey).(string)
	if !ok {
		e := errors.New("invalid username")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Service
	req.Username = username
	err := a.GetUser(req, &res)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

func (a *Auth) UpdateUserJSON(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserReq
	var res UpdateUserRes

	ctx := r.Context()
	username, ok := ctx.Value(userCtxKey).(string)
	if !ok {
		e := errors.New("invalid username")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Decode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Service
	req.Identifier.Username = username
	err = a.UpdateUser(req, &res)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

func (a *Auth) DeleteUserJSON(w http.ResponseWriter, r *http.Request) {
	var req DeleteUserReq
	var res DeleteUserRes

	ctx := r.Context()
	username, ok := ctx.Value(userCtxKey).(string)
	if !ok {
		e := errors.New("invalid username")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Service
	req.Identifier.Username = username
	err := a.DeleteUser(req, &res)
	if err != nil {
		e := errors.New("invalid username")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

// Account --------------------------------------------------------------------
func (a *Auth) CreateAccountJSON(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountReq
	var res CreateAccountRes

	// Decode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Service
	err = a.CreateAccount(req, &res)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

func (a *Auth) GetAccountsJSON(w http.ResponseWriter, r *http.Request) {
	var req GetAccountsReq
	var res GetAccountsRes

	// Service
	err := a.GetAccounts(req, &res)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

func (a *Auth) GetAccountJSON(w http.ResponseWriter, r *http.Request) {
	var req GetAccountReq
	var res GetAccountRes

	ctx := r.Context()
	slug, ok := ctx.Value(accountCtxKey).(string)
	if !ok {
		e := errors.New("invalid slug")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Service
	req.Slug = slug
	err := a.GetAccount(req, &res)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

func (a *Auth) UpdateAccountJSON(w http.ResponseWriter, r *http.Request) {
	var req UpdateAccountReq
	var res UpdateAccountRes

	ctx := r.Context()
	slug, ok := ctx.Value(accountCtxKey).(string)
	if !ok {
		e := errors.New("invalid slug")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Decode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Service
	req.Identifier.Slug = slug
	err = a.UpdateAccount(req, &res)
	if err != nil {
		a.Log().Error(err)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

func (a *Auth) DeleteAccountJSON(w http.ResponseWriter, r *http.Request) {
	var req DeleteAccountReq
	var res DeleteAccountRes

	ctx := r.Context()
	slug, ok := ctx.Value(accountCtxKey).(string)
	if !ok {
		e := errors.New("invalid slug")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Service
	req.Identifier.Slug = slug
	err := a.DeleteAccount(req, &res)
	if err != nil {
		e := errors.New("invalid slug")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

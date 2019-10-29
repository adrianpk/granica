package jsonrest

import (
	"encoding/json"
	"errors"
	"net/http"

	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
)

const (
	AccountCtxKey contextKey = "account"
)

func (ep *Endpoint) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req tp.CreateAccountReq
	var res tp.CreateAccountRes

	// Decode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Service
	err = ep.service.CreateAccount(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

func (ep *Endpoint) GetAccounts(w http.ResponseWriter, r *http.Request) {
	var req tp.GetAccountsReq
	var res tp.GetAccountsRes

	// Service
	err := ep.service.GetAccounts(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

func (ep *Endpoint) GetAccount(w http.ResponseWriter, r *http.Request) {
	var req tp.GetAccountReq
	var res tp.GetAccountRes

	ctx := r.Context()
	slug, ok := ctx.Value(AccountCtxKey).(string)
	if !ok {
		e := errors.New("invalid slug")
		ep.Log().Error(e)
		ep.writeResponse(w, res)
		return
	}

	// Service
	req.Slug = slug
	err := ep.service.GetAccount(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

func (ep *Endpoint) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var req tp.UpdateAccountReq
	var res tp.UpdateAccountRes

	ctx := r.Context()
	slug, ok := ctx.Value(AccountCtxKey).(string)
	if !ok {
		e := errors.New("invalid slug")
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
	req.Identifier.Slug = slug
	err = ep.service.UpdateAccount(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

func (ep *Endpoint) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	var req tp.DeleteAccountReq
	var res tp.DeleteAccountRes

	ctx := r.Context()
	slug, ok := ctx.Value(AccountCtxKey).(string)
	if !ok {
		e := errors.New("invalid slug")
		ep.Log().Error(e)
		ep.writeResponse(w, res)
		return
	}

	// Service
	req.Identifier.Slug = slug
	err := ep.service.DeleteAccount(req, &res)
	if err != nil {
		e := errors.New("invalid slug")
		ep.Log().Error(e)
		ep.writeResponse(w, res)
		return
	}

	// Output
	ep.writeResponse(w, res)
}

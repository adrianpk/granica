package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type (
	Response interface {
		GetErr() error
	}
)

const (
	createErr = "Cannot create entity"
	getAllErr = "Cannot get entity"
	getErr    = "Cannot get entity"
	updateErr = "Cannot update entity"
	deleteErr = "Cannot delete entity"
)

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
		e := errors.New("username not provided")
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
		e := errors.New("username not provided")
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
		e := errors.New("username not provided")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Service
	req.Identifier.Username = username
	err := a.DeleteUser(req, &res)
	if err != nil {
		e := errors.New("username not provided")
		a.Log().Error(e)
		a.writeResponse(w, res)
		return
	}

	// Output
	a.writeResponse(w, res)
}

// Output
func (a *Auth) writeResponse(w http.ResponseWriter, res interface{}) {
	// Marshalling
	o, err := a.toJSON(res)
	if err != nil {
		a.Log().Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(o)
}

func formatRequestBody(r *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return buf.String()
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	// If this is a POST, add post data
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

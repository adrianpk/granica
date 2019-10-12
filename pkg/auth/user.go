package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/mikrowezel/granica/internal/repo"
)

const (
	createErr = "Cannot create user"
	getAllErr = "Cannot get users"
	getErr    = "Cannot get user"
)

func (a *Auth) createUser(w http.ResponseWriter, r *http.Request) {
	// Unmarshal
	var uReq CreateUserReq
	err := json.NewDecoder(r.Body).Decode(&uReq)
	if err != nil {
		err = a.createUserResponse(w, r, nil, createErr, err)
		a.Log().Error(err)
		return
	}

	// Create a model
	u := uReq.toModel()

	// Repo
	repo, err := a.userRepo()
	if err != nil {
		err = a.createUserResponse(w, r, &u, createErr, err)
		a.Log().Error(err)
		return
	}

	err = repo.Create(&u)
	if err != nil {
		err = a.createUserResponse(w, r, &u, createErr, err)
		a.Log().Error(err)
		return
	}

	err = repo.Commit()
	if err != nil {
		err = a.createUserResponse(w, r, &u, createErr, err)
		a.Log().Error(err)
		return
	}

	// Output
	err = a.createUserResponse(w, r, &u, "", nil)
	if err != nil {
		a.Log().Error(err)
	}
}

func (a *Auth) getUsers(w http.ResponseWriter, r *http.Request) {
	// Repo
	repo, err := a.userRepo()
	if err != nil {
		err = a.getUsersResponse(w, r, nil, getAllErr, err)
		a.Log().Error(err)
		return
	}

	us, err := repo.GetAll()
	if err != nil {
		err = a.getUsersResponse(w, r, nil, getAllErr, err)
		a.Log().Error(err)
		return
	}

	err = repo.Commit()
	if err != nil {
		err = a.getUsersResponse(w, r, us, getAllErr, err)
		a.Log().Error(err)
		return
	}

	// Output
	err = a.getUsersResponse(w, r, us, "", nil)
	if err != nil {
		a.Log().Error(err)
	}
}

func (a *Auth) getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := ctx.Value(userCtxKey).(string)
	if !ok {
		e := errors.New("ID was not provided")
		err := a.getUserResponse(w, r, nil, getErr, e)
		a.Log().Error(err)
		return
	}

	// Repo
	repo, err := a.userRepo()
	if err != nil {
		err = a.getUserResponse(w, r, nil, getErr, err)
		a.Log().Error(err)
		return
	}

	u, err := repo.Get(id)
	if err != nil {
		err = a.getUserResponse(w, r, nil, getErr, err)
		a.Log().Error(err)
		return
	}

	err = repo.Commit()
	if err != nil {
		err = a.getUserResponse(w, r, nil, getErr, err)
		a.Log().Error(err)
		return
	}

	// Output
	err = a.getUserResponse(w, r, &u, "", nil)
	if err != nil {
		a.Log().Error(err)
	}
}

func (a *Auth) updateUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (a *Auth) deleteUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (a *Auth) errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	str := fmt.Sprintf(`{"data":"","error":"%s"}`, err.Error())
	fmt.Fprint(w, str)
	a.Log().Error(err, err.Error())
}

func (a *Auth) userRepo() (*repo.UserRepo, error) {
	rh, err := a.repoHandler()
	if err != nil {
		return nil, err
	}
	return rh.UserRepoNewTx()
}

func (a *Auth) repoHandler() (*repo.Repo, error) {
	h, ok := a.Handler("repo-handler")
	if !ok {
		return nil, errors.New("repo handler not available")
	}

	repo, ok := h.(*repo.Repo)
	if !ok {
		return nil, errors.New("invalid repo handler")
	}

	return repo, nil
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

package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/mikrowezel/granica/internal/model"
	"gitlab.com/mikrowezel/granica/internal/repo"
)

// createUser
/* using a JSON input like this:
{
  "username": "username",
  "password": "username@mail.com",
  "email": "username@mail.com",
  "emailConfirmation": "username@mail.com",
  "givenName": "name",
  "middleNames": "middle",
  "familyName": "family"
}*/
func (a *Auth) createUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	json.NewDecoder(r.Body).Decode(&u)

	// Get repo.
	repo, err := a.repoHandler()
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	a.Log().Info("Create user", "user", fmt.Sprintf("%+v", u))

	// Persist.
	err = repo.CreateUser(&u)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	err = repo.CommitTx()
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// Transform output.
	json, err := a.toJSON(&u)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// Output result.
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (a *Auth) getUsers(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	// Repo handler.
	//repo, err := a.repoHandler()
	//if err != nil {
	//a.errorResponse(w, r, err)
	//return
	//}

	// Transform output
	json, err := a.toJSON("")
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// Output result.
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (a *Auth) getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// User context
	_, ok := ctx.Value(userCtxKey).(string)
	if !ok {
		err := errors.New("no user context")
		a.errorResponse(w, r, err)
		return
	}

	// Repo handler.
	//repo, err := a.repoHandler()
	//if err != nil {
	//a.errorResponse(w, r, err)
	//return
	//}

	// Transform output
	json, err := a.toJSON("")
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// Output result.
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (a *Auth) updateUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (a *Auth) deleteUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (a *Auth) toJSON(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (a *Auth) errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	str := fmt.Sprintf(`{"data":"","error":"%s"}`, err.Error())
	fmt.Fprint(w, str)
	a.Log().Error(err, err.Error())
}

func (a *Auth) repoHandler() (*repo.Repo, error) {
	a.Log().Info("Handlers", "all", fmt.Sprintf("%+v", a.Handlers()))

	h, ok := a.Handler("repo-handler")
	if !ok {
		return nil, errors.New("repo handler not available")
	}

	a.Log().Info("Repo", "type", fmt.Sprintf("%+T", h))
	repo, ok := h.(*repo.Repo)
	if !ok {
		return nil, errors.New("invalid repo handler")
	}

	return repo, nil
}

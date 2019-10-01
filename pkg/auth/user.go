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
	// Unmarshal
	var ut CreateUserReq
	err := json.NewDecoder(r.Body).Decode(&ut)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// Create a model
	u := ut.toModel()

	// Persist
	repo, err := a.userRepo()
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	err = repo.Create(&u)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	err = repo.Commit()
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// Output result
	json, err := a.toJSON(&u)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

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

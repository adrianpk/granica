package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/mikrowezel/granica/internal/repo"
)

func (a *Auth) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get points from context.
	_, ok := ctx.Value(userCtxKey).(*AuthCtx)
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
	h, ok := a.Handler("repo-handler")
	if !ok {
		return nil, errors.New("Repo handler not available")
	}

	repo, ok := h.(*repo.Repo)
	if !ok {
		return nil, errors.New("invalidad repo handler")
	}

	return repo, nil
}

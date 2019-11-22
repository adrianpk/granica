package web

import (
	"net/http"

	//"github.com/davecgh/go-spew/spew"
	"gitlab.com/mikrowezel/backend/granica/internal/model"
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
	res := &tp.CreateUserRes{}
	res.FromModel(&model.User{}, "", nil)
	res.Action = ep.userCreateAction()

	// Wrap response
	wr := ep.OKRes(r, res)

	// Template
	ts, err := ep.TemplateFor(userRes, web.CreateTmpl)
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

func (ep *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req tp.CreateUserReq
	var res tp.CreateUserRes

	// TODO: Form data validation

	// Form to Req
	err := web.NewDecoder().Decode(&req.User, r.Form)
	//ep.Log().Debug("Resource", "values", spew.Sdump(res))

	// Req to model
	res.Action = ep.userCreateAction()

	// Template
	ts, err := ep.TemplateFor(userRes, web.CreateTmpl)
	if err != nil {
		ep.Redirect(w, r, "/")
		return
	}

	if err != nil {
		wr := ep.ErrRes(r, res, "Cannot create user")
		ep.Log().Error(err)

		// TODO: Use redirect instead.
		// Cleaner and avoids extra nesting level
		// Preserve form data using gorilla session
		err = ts.Execute(w, wr)
		if err != nil {
			ep.Log().Error(err)
			ep.Redirect(w, r, "/")
		}
		return
	}

	// Service
	err = ep.service.CreateUser(req, &res)
	if err != nil {
		wr := ep.ErrRes(r, res, "Cannot create user")
		ep.Log().Error(err)

		err = ts.Execute(w, wr)
		if err != nil {
			ep.Log().Error(err)
			ep.Redirect(w, r, "/")
		}
		return
	}

	// TODO: Flash message

	// Wrap response
	ep.Redirect(w, r, "/")

}

// GetUsers web endpoint.
func (ep *Endpoint) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req tp.GetUsersReq
	var res tp.GetUsersRes

	// Service
	err := ep.service.GetUsers(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.Redirect(w, r, "/")
		return
	}

	// Wrap response
	wr := ep.OKRes(r, res)

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

// Misc
// userCreateAction
func (ep *Endpoint) userCreateAction() web.Action {
	return web.Action{Target: UserPath(), Method: "POST"}
}

// userUpdateAction
func (ep *Endpoint) userUpdateAction(resource string, model web.Identifiable) web.Action {
	return web.Action{Target: UserPathSlug(model), Method: "PUT"}
}

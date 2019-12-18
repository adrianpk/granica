package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
	"gitlab.com/mikrowezel/backend/web"
)

const (
	userRes = "user"
)

const (
	UserCtxKey web.ContextKey = "user"
)

const (
	// Defined in 'assets/web/embed/i18n/xx.json'
	// Where xx is: de, en, es, pl.
	// TODO: Make it possible to use all locales.
	// Info
	UserCreatedInfoID = "user_created_info_msg"
	UserUpdatedInfoID = "user_updated_info_msg"
	UserDeletedInfoID = "user_deleted_info_msg"
	LoggedInInfoID    = "logged_in_info_msg"
	// Error
	CreateUserErrID  = "create_user_err_msg"
	IndexUsersErrID  = "get_all_users_err_msg"
	GetUserErrID     = "get_user_err_msg"
	UpdateUserErrID  = "update_user_err_msg"
	DeleteUserErrID  = "delete_user_err_msg"
	CredentialsErrID = "credentials_err_msg"
)

// IndexUsers web endpoint.
func (ep *Endpoint) IndexUsers(w http.ResponseWriter, r *http.Request) {
	var req tp.IndexUsersReq
	var res tp.IndexUsersRes

	// Service
	err := ep.service.IndexUsers(req, &res)
	if err != nil {
		ep.handleError(w, r, "/", IndexUsersErrID, err)
		return
	}

	// Wrap response
	wr := ep.OKRes(w, r, res, "")

	// Template
	ts, err := ep.TemplateFor(userRes, web.IndexTmpl)
	if err != nil {
		ep.handleError(w, r, "/", IndexUsersErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.handleError(w, r, "/", IndexUsersErrID, err)
		return
	}
}

func (ep *Endpoint) NewUser(w http.ResponseWriter, r *http.Request) {
	// Req & Res
	res := &tp.CreateUserRes{}
	res.IsNew = true
	res.Action = ep.userCreateAction()

	// Wrap response
	wr := ep.OKRes(w, r, res, "")

	// Template
	ts, err := ep.TemplateFor(userRes, web.NewTmpl)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}
}

func (ep *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req tp.CreateUserReq
	var res tp.CreateUserRes
	res.IsNew = true
	res.Action = ep.userCreateAction()

	// Input data to request struct
	err := ep.FormToModel(r, &req.User)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	// Service
	err = ep.service.CreateUser(req, &res)

	// Input validation errors
	if !res.Errors.IsEmpty() {
		ep.rerenderUserForm(w, r, res, web.NewTmpl)
		return
	}

	// Non validation errors
	if err != nil {
		ep.handleError(w, r, UserPath(), CreateUserErrID, err)
		return
	}

	m := ep.localize(r, UserCreatedInfoID)
	ep.RedirectWithFlash(w, r, UserPath(), m, web.InfoMT)
}

// ShowUser web endpoint.
func (ep *Endpoint) ShowUser(w http.ResponseWriter, r *http.Request) {
	var req tp.GetUserReq
	var res tp.GetUserRes

	// Identifier
	id, err := ep.getIdentifier(r)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	req = tp.GetUserReq{id}

	// Service
	err = ep.service.GetUser(req, &res)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Wrap response
	wr := ep.OKRes(w, r, res, "")

	// Template
	ts, err := ep.TemplateFor(userRes, web.ShowTmpl)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}
}

// EditUser web endpoint.
func (ep *Endpoint) EditUser(w http.ResponseWriter, r *http.Request) {
	// Req & Res
	var req tp.GetUserReq
	res := tp.GetUserRes{}

	// Identifier
	id, err := ep.getIdentifier(r)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	req = tp.GetUserReq{id}

	// Service
	err = ep.service.GetUser(req, &res)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Set additional values
	res.IsNew = false
	res.Action = ep.userUpdateAction(res)

	// Wrap response
	wr := ep.OKRes(w, r, res, "")

	// Template
	ts, err := ep.TemplateFor(userRes, web.EditTmpl)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}
}

// UpdateUser web endpoint.
func (ep *Endpoint) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req tp.UpdateUserReq
	var res tp.UpdateUserRes

	// Identifier
	id, err := ep.getIdentifier(r)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	req = tp.UpdateUserReq{id, tp.User{}}

	// Input data to request struct
	err = ep.FormToModel(r, &req.User)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	// Service
	err = ep.service.UpdateUser(req, &res)

	// Input validation errors
	if !res.Errors.IsEmpty() {
		ep.rerenderUserForm(w, r, res, web.NewTmpl)
		return
	}

	// Non validation errors
	if err != nil {
		ep.handleError(w, r, UserPath(), UpdateUserErrID, err)
		return
	}

	m := ep.localize(r, UserUpdatedInfoID)
	ep.RedirectWithFlash(w, r, UserPath(), m, web.InfoMT)
}

// InitDeleteUser web endpoint.
func (ep *Endpoint) InitDeleteUser(w http.ResponseWriter, r *http.Request) {
	var req tp.GetUserReq
	var res tp.GetUserRes

	// Identifier
	id, err := ep.getIdentifier(r)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	req = tp.GetUserReq{id}

	// Service
	err = ep.service.GetUser(req, &res)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Set additional values
	res.Action = ep.userDeleteAction(res)

	// Wrap response
	wr := ep.OKRes(w, r, res, "")

	// Template
	ts, err := ep.TemplateFor(userRes, web.InitDelTmpl)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}
}

// DeleteUser web endpoint.
func (ep *Endpoint) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req tp.DeleteUserReq
	var res tp.DeleteUserRes

	ctx := r.Context()
	slug, ok := ctx.Value(UserCtxKey).(string)
	if !ok {
		err := errors.New("no slug provided")
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	req = tp.DeleteUserReq{
		tp.Identifier{
			Slug: slug,
		},
	}

	// Service
	err := ep.service.DeleteUser(req, &res)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	m := ep.localize(r, UserDeletedInfoID)
	ep.RedirectWithFlash(w, r, UserPath(), m, web.InfoMT)
}

func (ep *Endpoint) InitSignInUser(w http.ResponseWriter, r *http.Request) {
	// Req & Res
	res := &tp.SigninUserRes{}
	res.Action = ep.userSigninAction()

	// Wrap response
	wr := ep.OKRes(w, r, res, "")

	// Template
	ts, err := ep.TemplateFor(userRes, web.SigninTmpl)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}
}

// SignInUser web endpoint.
func (ep *Endpoint) SignInUser(w http.ResponseWriter, r *http.Request) {
	var req tp.SigninUserReq
	var res tp.SigninUserRes

	// Input data to request struct
	err := ep.FormToModel(r, &req.Signin)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	// Service
	err = ep.service.SignInUser(req, &res)
	if err != nil {
		ep.handleError(w, r, UserPathSignin(), CredentialsErrID, err)
		return
	}

	m := ep.localize(r, LoggedInInfoID)
	ep.RedirectWithFlash(w, r, UserPath(), m, web.InfoMT)
}

func (ep *Endpoint) rerenderUserForm(w http.ResponseWriter, r *http.Request, res interface{}, template string) {
	wr := ep.ErrRes(w, r, res, InputValuesErrID, nil)

	ts, err := ep.TemplateFor(userRes, template)
	if err != nil {
		ep.handleError(w, r, UserPath(), InputValuesErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	return
}

// Localization - I18N
func (ep *Endpoint) localize(r *http.Request, msgID string) string {
	l := ep.Localizer(r)
	if l == nil {
		ep.Log().Warn("No localizer available")
		return msgID
	}

	t, _, err := l.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID: msgID,
	})

	if err != nil {
		ep.Log().Error(err)
		return msgID
	}

	//ep.Log().Debug("Localized message", "value", t, "lang", lang)

	return t
}

func (ep *Endpoint) localizeMessageID(l *i18n.Localizer, messageID string) (string, error) {
	return l.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

// Misc
func (ep *Endpoint) getIdentifier(r *http.Request) (identifier tp.Identifier, err error) {
	ctx := r.Context()
	slug, ok := ctx.Value(UserCtxKey).(string)
	if !ok {
		err := errors.New("no slug provided")
		return tp.Identifier{}, err
	}

	return tp.Identifier{
		Slug: slug,
	}, nil
}

// userCreateAction
func (ep *Endpoint) userCreateAction() web.Action {
	return web.Action{Target: fmt.Sprintf("%s", UserPath()), Method: "POST"}
}

// userUpdateAction
func (ep *Endpoint) userUpdateAction(model web.Identifiable) web.Action {
	return web.Action{Target: UserPathSlug(model), Method: "PUT"}
}

// userDeleteAction
func (ep *Endpoint) userDeleteAction(model web.Identifiable) web.Action {
	return web.Action{Target: UserPathSlug(model), Method: "DELETE"}
}

// userSignupAction
func (ep *Endpoint) userSignupAction() web.Action {
	return web.Action{Target: UserPathSignup(), Method: "POST"}
}

// userSigninAction
func (ep *Endpoint) userSigninAction() web.Action {
	return web.Action{Target: UserPathSignin(), Method: "POST"}
}

func (ep *Endpoint) handleError(w http.ResponseWriter, r *http.Request, redirPath, msgID string, err error) {
	m := ep.localize(r, msgID)
	ep.RedirectWithFlash(w, r, redirPath, m, web.ErrorMT)
	ep.Log().Error(err)
}

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
	// Error
	CannotProcErrID = "cannot_proc_err_msg"
	CreateUserErrID = "create_user_err_msg"
	IndexUsersErrID = "get_all_users_err_msg"
	GetUserErrID    = "get_user_err_msg"
	UpdateUserErrID = "update_user_err_msg"
	DeleteUserErrID = "delete_user_err_msg"
)

func (ep *Endpoint) NewUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve stored form data if exists
	// It avoids filling in the form again after submissions errors.
	userForm := ep.RestoreUserForm(r, web.CreateUserStoreKey)

	// Req & Res
	res := &tp.CreateUserRes{}
	res.FromTransport(&userForm, "", nil)
	res.IsNew = true
	res.Action = ep.userCreateAction()

	// Wrap response
	wr := ep.OKRes(r, res, "")

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

	// TODO: Form data validation

	// Form to Req
	err := ep.FormToModel(r, &req.User)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	//ep.Log().Info(spew.Sdump(&req.User))

	// Store form data
	ep.StoreUserForm(r, w, web.CreateUserStoreKey, req.User)

	// Service
	err = ep.service.CreateUser(req, &res)
	if err != nil {
		ep.handleError(w, r, UserPathNew(), CreateUserErrID, err)
		return
	}

	// Operation succeded: form data can be cleared.
	ep.ClearUserForm(r, w, web.CreateUserStoreKey)

	m := ep.localize(r, UserCreatedInfoID)
	ep.RedirectWithFlash(w, r, UserPath(), m, web.InfoMT)
}

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
	wr := ep.OKRes(r, res, "")

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

// EditUser web endpoint.
func (ep *Endpoint) EditUser(w http.ResponseWriter, r *http.Request) {
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

	// Set action
	res.Action = ep.userUpdateAction(userRes, res)

	// Wrap response
	wr := ep.OKRes(r, res, "")

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
	wr := ep.OKRes(r, res, "")

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

// UpdateUser web endpoint.
func (ep *Endpoint) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ep.Log().Info("Updating user...")
	var req tp.UpdateUserReq
	var res tp.UpdateUserRes

	// TODO: Form data validation

	// Store form data if exists
	ep.StoreUserForm(r, w, web.UpdateUserStoreKey, req.User)

	// Identifier
	id, err := ep.getIdentifier(r)
	if err != nil {
		ep.handleError(w, r, UserPath(), GetUserErrID, err)
		return
	}

	req = tp.UpdateUserReq{id, tp.User{}}

	// Form to Req
	err = ep.FormToModel(r, &req.User)
	if err != nil {
		ep.handleError(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	// Service
	err = ep.service.UpdateUser(req, &res)
	if err != nil {
		ep.handleError(w, r, UserPathNew(), CreateUserErrID, err)
		return
	}

	// Operation succeded: form data can be cleared.
	ep.ClearUserForm(r, w, web.UpdateUserStoreKey)

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

	// Set action
	res.Action = ep.userDeleteAction(userRes, res)

	// Wrap response
	wr := ep.OKRes(r, res, "")

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

// Localization - I18N

func (ep *Endpoint) localize(r *http.Request, msgID string) string {
	l := ep.Localizer(r)
	if l == nil {
		ep.Log().Warn("No localizer available")
		return msgID
	}

	t, lang, err := l.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID: msgID,
	})

	if err != nil {
		ep.Log().Error(err)
		return msgID
	}

	ep.Log().Debug("Localized message", "value", t, "lang", lang)

	return t
}

func (ep *Endpoint) localizeMessageID(l *i18n.Localizer, messageID string) (string, error) {
	return l.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

// Form data session helpers
func (ep *Endpoint) StoreUserForm(r *http.Request, w http.ResponseWriter, key string, userForm tp.User) (ok bool) {
	s := ep.GetSession(r)
	s.Values[key] = userForm
	err := s.Save(r, w)
	if err != nil {
		return true
	}
	return false
}

func (ep *Endpoint) RestoreUserForm(r *http.Request, key string) tp.User {
	s := ep.GetSession(r)
	v := s.Values[key]

	user, ok := v.(tp.User)
	if ok {
		//ep.Log().Debug("Stored form data", "value", spew.Sdump(user))
		return user
	}

	ep.Log().Info("No stored form data", "key", key)
	return tp.User{}
}

func (ep *Endpoint) ClearUserForm(r *http.Request, w http.ResponseWriter, key string) (ok bool) {
	s := ep.GetSession(r)
	delete(s.Values, key)
	err := s.Save(r, w)
	if err != nil {
		return true
	}
	return false
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
func (ep *Endpoint) userUpdateAction(resource string, model web.Identifiable) web.Action {
	return web.Action{Target: UserPathSlug(model), Method: "PUT"}
}

// userDeleteAction
func (ep *Endpoint) userDeleteAction(resource string, model web.Identifiable) web.Action {
	return web.Action{Target: UserPathSlug(model), Method: "DELETE"}
}

func (ep *Endpoint) handleError(w http.ResponseWriter, r *http.Request, redirPath, msgID string, err error) {
	m := ep.localize(r, msgID)
	ep.RedirectWithFlash(w, r, redirPath, m, web.ErrorMT)
	ep.Log().Error(err)
}

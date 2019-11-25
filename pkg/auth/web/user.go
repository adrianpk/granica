package web

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
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
	// TODO: Make it possible to use all available locales.
	CannotProcErrID  = "cannot_proc_err_msg"
	CreateUserErrID  = "create_user_err_msg"
	GetAllUsersErrID = "get_all_users_err_msg"
	GetUserErrID     = "get_user_err_msg"
	UpdateUserErrID  = "update_user_err_msg"
	DeleteUserErrID  = "delete_user_err_msg"
)

func (ep *Endpoint) InitCreateUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve stored form data if exists
	// It avoids filling in the form again after submissions errors.
	userForm := ep.RestoreUserForm(r, web.CreateUserStoreKey)

	// Retrieve flash data if exists
	f := ep.RestoreFlash(r)
	if !f.IsEmpty() {
		// Just only logging at the moment
		ep.Log().Debug("Session flash data", "value", spew.Sdump(f))
	}

	// Req & Res
	res := &tp.CreateUserRes{}
	res.FromTransport(&userForm, "", nil)
	res.Action = ep.userCreateAction()

	// Wrap response
	wr := ep.OKRes(r, res)

	// Template
	ts, err := ep.TemplateFor(userRes, web.CreateTmpl)
	if err != nil {
		m := ep.localizeMsg(r, CannotProcErrID)
		ep.StoreFlash(r, w, m, web.ErrorMT)
		ep.Redirect(w, r, UserPath())
		ep.Log().Error(err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		m := ep.localizeMsg(r, CannotProcErrID)
		ep.StoreFlash(r, w, m, web.ErrorMT)
		ep.Redirect(w, r, UserPath())
		ep.Log().Error(err)
		return
	}
}

func (ep *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req tp.CreateUserReq
	var res tp.CreateUserRes

	// TODO: Form data validation

	// Store form data if exists
	// It avoids filling in the form again amter submissions errors.
	ep.StoreUserForm(r, w, web.CreateUserStoreKey, req.User)

	// Form to Req
	err := web.FormToModel(r, &req.User)
	if err != nil {
		ep.Redirect(w, r, UserPathNew())
		return
	}

	// Service
	err = ep.service.CreateUser(req, &res)
	if err != nil {
		m := ep.localizeMsg(r, CreateUserErrID)
		ep.StoreFlash(r, w, m, web.ErrorMT)
		ep.Redirect(w, r, UserPathNew())
		ep.Log().Error(err)
		return
	}

	// Operation succed: form data can be cleared.
	ep.ClearUserForm(r, w, web.CreateUserStoreKey)

	// Flash message
	ep.StoreFlash(r, w, "Sample message", web.InfoMT)

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

// Localization - I18N
func (ep *Endpoint) localizeMsg(r *http.Request, msgID string) string {
	l, ok := web.GetI18NLocalizer(r)
	if !ok {
		// FIX: Do something: Return default message?
		ep.Log().Warn("I18N localizer not available")
	}

	t, lang, err := l.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID: msgID,
	})

	if err != nil {
		ep.Log().Error(err)
	}

	ep.Log().Debug("Localized message", "value", t, "lang", lang)

	return t
}

func (ep *Endpoint) localizeMessageID(l *i18n.Localizer, messageID string) (string, error) {
	return l.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

// Flash
func (ep *Endpoint) StoreFlash(r *http.Request, w http.ResponseWriter, message string, mt web.MsgType) (ok bool) {
	s := ep.GetSession(r)

	s.Values[web.FlashStoreKey] = ep.MakeFlash(message, mt)
	err := s.Save(r, w)
	if err != nil {
		ep.Log().Error(err)
		return true
	}

	return false
}

func (ep *Endpoint) RestoreFlash(r *http.Request) web.FlashData {
	s := ep.GetSession(r)
	v := s.Values[web.FlashStoreKey]

	user, ok := v.(web.FlashData)
	if ok {
		ep.Log().Debug("Stored flash", "value", spew.Sdump(user))
		return web.FlashData{}
	}

	ep.Log().Info("No stored flash", "key", web.FlashStoreKey)
	return web.FlashData{}
}

func (ep *Endpoint) ClearFlash(r *http.Request, w http.ResponseWriter, message string, mt web.MsgType) (ok bool) {
	s := ep.GetSession(r)
	delete(s.Values, web.FlashStoreKey)
	err := s.Save(r, w)
	if err != nil {
		return true
	}
	return false
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
		ep.Log().Debug("Stored form data", "value", spew.Sdump(user))
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
// userCreateAction
func (ep *Endpoint) userCreateAction() web.Action {
	return web.Action{Target: fmt.Sprintf("%s", UserPath()), Method: "POST"}
}

// userUpdateAction
func (ep *Endpoint) userUpdateAction(resource string, model web.Identifiable) web.Action {
	return web.Action{Target: UserPathSlug(model), Method: "PUT"}
}

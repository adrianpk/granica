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
	createUserKey = "create-user"
	updateUserKey = "update-user"
)

func (ep *Endpoint) InitCreateUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve stored form data if exists
	// It avoids filling in the form again after submissions errors.
	userFD := ep.restoreUserFD(r, createUserKey)

	// Req & Res
	res := &tp.CreateUserRes{}
	res.FromTransport(&userFD, "", nil)
	res.Action = ep.userCreateAction()

	// Wrap response
	wr := ep.OKRes(r, res)

	// Template
	ts, err := ep.TemplateFor(userRes, web.CreateTmpl)
	if err != nil {
		m := ep.errMsg(r, web.CannotProcErr, userRes)
		ep.storeErrorFlash(r, w, m)
		ep.Redirect(w, r, UserPath())
		ep.Log().Error(err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		m := ep.errMsg(r, web.CannotProcErr, userRes)
		ep.storeErrorFlash(r, w, m)
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
	// It avoids filling in the form again after submissions errors.
	ep.storeUserFD(r, w, createUserKey, req.User)

	// Form to Req
	err := web.FormToModel(r, &req.User)
	if err != nil {
		ep.Redirect(w, r, UserPathNew())
		return
	}

	// Service
	err = ep.service.CreateUser(req, &res)
	if err != nil {
		m := ep.errMsg(r, web.CreateErr, userRes)
		ep.storeErrorFlash(r, w, m)
		ep.Redirect(w, r, UserPathNew())
		ep.Log().Error(err)
		return
	}

	ep.clearUserFD(r, w, createUserKey)
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
	return web.Action{Target: fmt.Sprintf("%s", UserPath()), Method: "POST"}
}

// userUpdateAction
func (ep *Endpoint) userUpdateAction(resource string, model web.Identifiable) web.Action {
	return web.Action{Target: UserPathSlug(model), Method: "PUT"}
}

// Flash messages
func (ep *Endpoint) errMsg(r *http.Request, errType, resource string) string {
	l, ok := web.GetI18NLocalizer(r)
	if !ok {
		// FIX: Do something: Return default message?
		ep.Log().Warn("I18N localizer not available")
	}

	// Message
	id := fmt.Sprintf("%s_err_msg", errType)

	t, lang, err := l.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID: id,
		//TemplateData: map[string]string{
		//"Name": resource,
		//},
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
func (ep *Endpoint) storeErrorFlash(r *http.Request, w http.ResponseWriter, message string) (ok bool) {
	ep.Log().Debug("*Endpoint.storeErrorFlash not implemented")
	return false
}

// Form data session helpers
func (ep *Endpoint) storeUserFD(r *http.Request, w http.ResponseWriter, key string, userFD tp.User) (ok bool) {
	s := ep.GetSession(r)
	s.Values[key] = userFD
	err := s.Save(r, w)
	if err != nil {
		return true
	}
	return false
}

func (ep *Endpoint) restoreUserFD(r *http.Request, key string) tp.User {
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

func (ep *Endpoint) clearUserFD(r *http.Request, w http.ResponseWriter, key string) (ok bool) {
	s := ep.GetSession(r)
	delete(s.Values, key)
	err := s.Save(r, w)
	if err != nil {
		return true
	}
	return false
}

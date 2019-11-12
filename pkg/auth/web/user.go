package web

import (
	"html/template"
	"net/http"

	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
)

const (
	UserCtxKey contextKey = "user"
)

// TODO: This is a work in progress and the implementation is still unclean.
// - Common logic to all handlers will be extracted and generalized.
// - Templates will be embedded (https://github.com/markbates/pkger)
// - Error condition will load a flash message and render/redirec page as appropiate.
// - Templates will be beautified using tailwind classes.
// - To consider: allow loading of templates from external filepathso
// 	* External / embedded configurable by envar.
// - Finally, after a pattern emerges, all resources needed\
// to generate endpoint handlers and templates will be automated
// using mw-cli: https://gitlab.com/mikrowezel/backend/cli
func (ep *Endpoint) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req tp.GetUsersReq
	var res tp.GetUsersRes

	// Service
	err := ep.service.GetUsers(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, err.Error()) // FIX: Implement a redirect.
		return
	}

	//// Output
	//ep.writeResponse(w, res)

	// Template paths
	files := []string{
		"./assets/web/resource/base.layout.tmpl",
		"./assets/web/resource/flash.partial.tmpl",
		"./assets/web/resource/user/index.page.tmpl",
		"./assets/web/resource/user/header.partial.tmpl",
		"./assets/web/resource/user/ctxbar.partial.tmpl",
		"./assets/web/resource/user/list.partial.tmpl",
	}

	// Parse templates
	ts, err := template.New("base.layout.tmpl").Funcs(pathFxs).ParseFiles(files...)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, err.Error()) // FIX: Implement a redirect.
		return
	}

	//ep.Log().Info("Response object", "val", spew.Sdump(res))

	// Execute templates
	err = ts.Execute(w, res)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, err.Error()) // FIX: Implement a redirect.
		return
	}
}

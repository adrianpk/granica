package web

import (
	"html/template"
	"net/http"

	tp "gitlab.com/mikrowezel/backend/granica/pkg/auth/transport"
)

type (
	// WRes stands for wrapped response
	WRes struct {
		Data  interface{}
		Flash []FlashData
		Err   error
	}

	// Flash data to present in page
	FlashData struct {
		Type MsgType
		Msg  string
	}

	// MsgType stands for message type
	MsgType string
)

const (
	InfoMT  MsgType = "info"
	WarnMT  MsgType = "warn"
	ErrorMT MsgType = "error"
	DebugMT MsgType = "debug"
)

const (
	UserCtxKey contextKey = "user"
)

// TODO: This is a work in progress and the implementation is still unclean.
// - Common logic to all handlers will be extracted and generalized.
// 	* Move to a mw library or same login in project using cli generator?
// - Templates w̶i̶l̶l̶ ̶b̶e̶  are now embedded (https://github.com/markbates/pkger)
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
		"./assets/web/resource/layout/base.tmpl",
		"./assets/web/resource/partial/_flash.tmpl",
		"./assets/web/resource/user/index.tmpl",
		"./assets/web/resource/user/_header.tmpl",
		"./assets/web/resource/user/_ctxbar.tmpl",
		"./assets/web/resource/user/_list.tmpl",
	}

	// Parse templates
	ts, err := template.New("base.tmpl").Funcs(pathFxs).ParseFiles(files...)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, err.Error()) // FIX: Implement a redirect.
		return
	}

	//ep.Log().Info("Response object", "val", spew.Sdump(res))

	// Execute templates
	wr := ep.okRes(res)
	err = ts.Execute(w, wr)
	if err != nil {
		ep.Log().Error(err)
		ep.writeResponse(w, err.Error()) // FIX: Implement a redirect.
		return
	}
}

func (ep *Endpoint) okRes(data interface{}, msgs ...string) WRes {
	fls := []FlashData{}
	for _, m := range msgs {
		fls = append(fls, ep.makeFlash(m, InfoMT))
	}

	return WRes{
		Data:  data,
		Flash: fls,
		Err:   nil,
	}
}

func (ep *Endpoint) wrap(data interface{}, msg string, msgType MsgType, err error) WRes {
	return WRes{
		Data:  data,
		Flash: []FlashData{ep.makeFlash(msg, msgType)},
		Err:   err,
	}
}

func (ep *Endpoint) makeFlash(msg string, msgType MsgType) FlashData {
	return FlashData{
		Type: msgType,
		Msg:  msg,
	}
}

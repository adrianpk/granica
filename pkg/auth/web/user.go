package web

import (
	"errors"
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
	userRes = "user"
)

const (
	UserCtxKey contextKey = "user"
)

// TODO: This is a work in progress and the implementation is still unclean.
// - Common logic to all handlers will be extracted and generalized.
// 	* Move to a̶ ̶m̶w̶ ̶l̶i̶b̶r̶a̶r̶y̶ ̶o̶r̶  generate code in project using cli generator?
//    At the moment I prefer an approach that limits the amount of dependenciesi
//    even if it increases LOC of the project.
// - Templates w̶i̶l̶l̶ ̶b̶e̶  are now embedded (https://github.com/markbates/pkger)
// - Error condition will load a flash message and render/redirec page as appropiate.
// - T̶e̶m̶p̶l̶a̶t̶e̶s̶ ̶w̶i̶l̶l̶ ̶b̶e̶ ̶b̶e̶a̶u̶t̶i̶f̶i̶e̶d̶ ̶u̶s̶i̶n̶g̶ ̶t̶a̶i̶l̶w̶i̶n̶d̶ ̶c̶l̶a̶s̶s̶e̶s̶.
//  * Done
// - To consider: allow loading of templates from external filepath.
// T̶o̶ ̶c̶o̶n̶s̶i̶d̶e̶r̶:̶ ̶a̶l̶l̶o̶w̶ ̶l̶o̶a̶d̶i̶n̶g̶ ̶o̶f̶ ̶t̶e̶m̶p̶l̶a̶t̶e̶s̶ ̶f̶r̶o̶m̶ ̶e̶x̶t̶e̶r̶n̶a̶l̶ ̶f̶i̶l̶e̶p̶a̶t̶h̶.̶
// 	* E̶x̶t̶e̶r̶n̶a̶l̶ ̶/̶ ̶e̶m̶b̶e̶d̶d̶e̶d̶ ̶c̶o̶n̶f̶i̶g̶u̶r̶a̶b̶l̶e̶ ̶b̶y̶ ̶e̶n̶v̶a̶r̶.̶
//  * Discarded for now.
// - Finally, after a pattern emerges, all resources needed
// 	 to generate endpoint handlers and templates will be automated
//   using mw-cli: https://gitlab.com/mikrowezel/backend/cli
func (ep *Endpoint) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req tp.GetUsersReq
	var res tp.GetUsersRes

	// Service
	err := ep.service.GetUsers(req, &res)
	if err != nil {
		ep.Log().Error(err)
		ep.redirect(w, r, "/")
		return
	}

	// Wrap response
	wr := ep.okRes(res)

	// Template
	ts, err := ep.templateFor(userRes, indexTmpl)
	if err != nil {
		ep.redirect(w, r, "/")
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.Log().Error(err)
		ep.redirect(w, r, "/")
	}
}

func (ep *Endpoint) templateFor(res, name string) (*template.Template, error) {
	key := ep.template(userRes, indexTmpl)

	t, ok := ep.templates[key]
	if !ok {
		err := errors.New("canot get template")
		ep.Log().Error(err, "resource", res, "template", name)
		return nil, err
	}

	return t, nil
}

// okRes builds an OK response including data and cero, one  or more messages.
// All messages are assumed of type info therefore flashes will be also of this type.
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

// multiRes builds a multiType response including data and cero, one  or more messages.
// Generated flash messages will be created according to the values passed in the parameter map.
// i.e.: map[string]MsgType{"Action processed": InfoMT, "Remember to update profile": WarnMT}
func (ep *Endpoint) multiRes(data interface{}, msgs map[string]MsgType) WRes {
	fls := []FlashData{}
	for m, t := range msgs {
		fls = append(fls, ep.makeFlash(m, t))
	}

	return WRes{
		Data:  data,
		Flash: fls,
		Err:   nil,
	}
}

// errRes builds an error response including data and cero, one  or more messages.
// All messages are assumed of type error therefore flashes will be also of this type.
func (ep *Endpoint) errRes(data interface{}, msgs ...string) WRes {
	fls := []FlashData{}
	for _, m := range msgs {
		fls = append(fls, ep.makeFlash(m, ErrorMT))
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

// Redirect to url.
func (ep *Endpoint) redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 302)
}

package web

import (
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
		ep.writeResponse(w, err.Error()) // FIX: Implement a redirect.
		return
	}

	//// Output
	//ep.writeResponse(w, res)

	// Execute templates
	wr := ep.okRes(res)

	tmplKey := ep.makeTmplKey(userRes, indexTmpl)
	ts, ok := ep.parsed[tmplKey]
	if !ok {
		// do something
	}

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

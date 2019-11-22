package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/markbates/pkger"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/log"
	"golang.org/x/text/message"
)

type (
	Endpoint struct {
		ctx        context.Context
		cfg        *config.Config
		log        *log.Logger
		templates  TemplateSet
		templateFx template.FuncMap
		store      *sessions.CookieStore
		storeKey   string
	}

	TemplateSet    map[string]*template.Template
	TemplateGroups map[string]map[string][]string
)

type Action struct {
	Target string
	Method string
}

type (
	ContextKey string
)

type (
	// WrappedRes stands for wrapped response
	WrappedRes struct {
		Data  interface{}
		Flash []FlashData
		CSRF  map[string]interface{}
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

type (
	Identifiable interface {
		GetSlug() string
	}
)

const (
	GetMethod    = "GET"
	PostMethod   = "POST"
	PutMethod    = "PUT"
	PatchMethod  = "PATCH"
	DeleteMethod = "DELETE"
)

const (
	templateDir = "/assets/web/embed/template"
	layoutDir   = "layout"
	layoutKey   = "layout"
	pageKey     = "page"
	partialKey  = "partial"
)

const (
	IndexTmpl   = "index.tmpl"
	CreateTmpl  = "create.tmpl"
	UpdateTmpl  = "update.tmpl"
	ShowTmpl    = "show.tmpl"
	InitDelTmpl = "initdel.tmpl"
)

const (
	InfoMT  MsgType = "info"
	WarnMT  MsgType = "warn"
	ErrorMT MsgType = "error"
	DebugMT MsgType = "debug"
)

const (
	I18NorCtxKey ContextKey = "i18n"
)

const (
	GetErrFmt    = "Cannot create %s."
	GetAllErrFmt = "Cannot get list of %s."
	CreateErrFmt = "Cannot get %s."
	UpdateErrFmt = "Cannot update %s."
	DeleteErrFmt = "Cannot delete %s."
)

func MakeEndpoint(ctx context.Context, cfg *config.Config, log *log.Logger, templateFx template.FuncMap) (*Endpoint, error) {
	ep := Endpoint{
		ctx:        ctx,
		cfg:        cfg,
		log:        log,
		templateFx: templateFx,
	}

	// Cookie store
	ep.makeCookieStore()

	// Load
	ts, err := ep.loadTemplates()
	if err != nil {
		return &ep, err
	}

	// Classify
	tg := ep.classifyTemplates(ts)

	// Parse
	ep.parseTemplates(ts, tg)

	return &ep, nil
}

func (ep *Endpoint) Ctx() context.Context {
	return ep.ctx
}

func (ep *Endpoint) Cfg() *config.Config {
	return ep.cfg
}

func (ep *Endpoint) Log() *log.Logger {
	return ep.log
}

func (ep *Endpoint) Templates() TemplateSet {
	return ep.templates
}

func (ep *Endpoint) TemplatesFx() template.FuncMap {
	return ep.templateFx
}

func (ep *Endpoint) Store() *sessions.CookieStore {
	return ep.store
}

func (wr *WrappedRes) addCSRF(r *http.Request) {
	wr.CSRF = map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}
}

// loadTemplates from embedded filesystem (pkger)
// under '/assets/web/embed/template'
func (ep *Endpoint) loadTemplates() (TemplateSet, error) {
	tmpls := make(TemplateSet)

	err := pkger.Walk(templateDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				ep.Log().Error(err, "msg", "Cannot load template", "path", path)
				return err
			}

			if filepath.Ext(path) == ".tmpl" {
				list := filepath.SplitList(path)
				base := fmt.Sprintf("%s:%s", list[0], templateDir)
				p, _ := filepath.Rel(base, path)

				ep.Log().Info("Reading template", "path", p)

				tmpls[p] = template.New(base)
				return nil
			}

			//e.Log().Warn("Not a valid template", "path", path)

			return nil
		})

	if err != nil {
		ep.Log().Error(err, "msg", "Cannot load templates", "path")
		return tmpls, err
	}

	return tmpls, nil
}

// classifyTemplates grouping them,
// first by type (layout, partial and page)
// and then by resource.
func (ep *Endpoint) classifyTemplates(ts TemplateSet) TemplateGroups {
	all := make(TemplateGroups)
	last := ""
	keys := ep.tmplsKeys(ts)

	for _, path := range keys {
		p := "./assets/web/embed/template" + "/" + path

		//e.Log().Debug("Classifying", "path", path)

		fileDir := filepath.Dir(path)
		fileName := filepath.Base(path)

		if fileDir != last {

			if _, ok := all[fileDir]; !ok {
				all[fileDir] = make(map[string][]string)
			}

			if isValidTemplateFile(path) {
				if isPartial(fileName) {
					all[fileDir][partialKey] = append(all[fileDir][partialKey], p)

				} else if isLayout(fileDir) {
					all[layoutDir][layoutKey] = append(all[layoutDir][layoutKey], p)

				} else {
					all[fileDir][pageKey] = append(all[fileDir][pageKey], p)
				}
			}
		}
	}

	return all
}

func (ep *Endpoint) tmplsKeys(ts TemplateSet) []string {
	keys := make([]string, 0, len(ts))
	for k, _ := range ts {
		keys = append(keys, k)
	}
	return keys
}

// parseTemplates parses template sets for each resource.
func (ep *Endpoint) parseTemplates(ts TemplateSet, tg TemplateGroups) {
	ep.templates = make(TemplateSet)
	layout := tg[layoutDir][layoutKey][0]

	for k, ts := range tg {
		pages := ts[pageKey]
		partials := ts[partialKey]

		for _, t := range pages {
			if k != layoutDir {
				ep.parseTemplate(t, partials, layout, ep.templateFx)
			}
		}
	}
}

func (ep *Endpoint) parseTemplate(page string, partials []string, layout string, funcs template.FuncMap) {
	parse := "base.tmpl"
	all := make([]string, 10)
	all = append(all, page)
	all = append(all, partials...)
	all = append(all, layout)
	trimSlice(&all)

	t, err := template.New(parse).Funcs(funcs).ParseFiles(all...)
	if err != nil {
		ep.Log().Error(err, "Error parsing template set", "page", page)
	}

	base := fmt.Sprintf(".%s", templateDir)
	p, _ := filepath.Rel(base, page)

	ep.Log().Info("Parsed template set", "path", p)

	ep.templates[page] = t
}

func trimSlice(slice *[]string) {
	newSlice := make([]string, 0, len(*slice))
	for _, val := range *slice {
		switch val {
		case "":
		default:
			newSlice = append(newSlice, val)
		}
	}
	*slice = newSlice
}

func isValidTemplateFile(fileName string) bool {
	return strings.HasSuffix(fileName, ".tmpl") && !strings.HasPrefix(fileName, ".")
}

func isPartial(fileName string) bool {
	return strings.HasPrefix(fileName, "_")
}

func isLayout(fileDir string) bool {
	return strings.HasPrefix(fileDir, "layout")
}

// Output
func (ep *Endpoint) writeResponse(w http.ResponseWriter, res interface{}) {
	// Marshalling
	o, err := ep.toJSON(res)
	if err != nil {
		ep.Log().Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(o)
}

func (ep *Endpoint) toJSON(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func formatRequestBody(r *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	return buf.String()
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	// If this is a POST, add post data
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

// Cookie store
func (ep *Endpoint) makeCookieStore() {
	k := ep.Cfg().ValOrDef("web.cookiestore.key", "")
	if k == "" {
		k = ep.genAES256Key()
		ep.Log().Debug("New cookie store random key", "value", k)
		csEnVar := fmt.Sprintf("%s_COOKIESTORE_KEY", "GRN")
		ep.Log().Info("Set a custom cookie store key using a 32 char string stored as an envar", "envvar", csEnVar)
	}

	ep.storeKey = k
	ep.Log().Debug("Cookie store key", "value", k)
	sessions.NewCookieStore([]byte(k))
}

func (ep *Endpoint) genAES256Key() string {
	const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		lenght    = 32
		indexBits = 6                // 6 bits to represent a letter index
		indexMask = 1<<indexBits - 1 // All 1-bits, as many as letterIdxBits
		indexMax  = 63 / indexBits   // # of letter indices fitting in 63 bits
	)
	src := rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(32)
	for i, cache, remain := lenght-1, src.Int63(), indexMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), indexMax
		}
		if idx := int(cache & indexMask); idx < len(allowed) {
			sb.WriteByte(allowed[idx])
			i--
		}
		cache >>= indexBits
		remain--
	}

	return sb.String()
}

// Templates

func (ep *Endpoint) TemplateFor(res, name string) (*template.Template, error) {
	key := ep.template(res, name)

	t, ok := ep.templates[key]
	if !ok {
		err := errors.New("canot get template")
		ep.Log().Error(err, "resource", res, "template", name)
		return nil, err
	}

	return t, nil
}

func (ep *Endpoint) template(resource, template string) (tmplKey string) {
	return fmt.Sprintf(".%s/%s/%s", templateDir, resource, template)
}

// I18N
func I18NGetAllErrMsg(r *http.Request, resource string) string {
	return I18NErrMsg(r, resource, GetAllErrFmt)
}

func I18NGetErrMsg(r *http.Request, resource string) string {
	return I18NErrMsg(r, resource, GetErrFmt)
}

func I18NCreateErrMsg(r *http.Request, resource string) string {
	return I18NErrMsg(r, resource, CreateErrFmt)
}

func I18NUpdateErrMsg(r *http.Request, resource string) string {
	return I18NErrMsg(r, resource, UpdateErrFmt)
}

func I18NDeleteErrMsg(r *http.Request, resource string) string {
	return I18NErrMsg(r, resource, DeleteErrFmt)
}

func I18NErrMsg(r *http.Request, resource, errFmt string) string {
	return fmt.Sprintf(errFmt, resource)
	mp, ok := GetI18Nor(r)
	if ok {
		return mp.Sprintf(errFmt, resource)
	}

	return fmt.Sprintf(errFmt, resource)
}

func GetI18Nor(r *http.Request) (mp *message.Printer, ok bool) {
	mp, ok = r.Context().Value(I18NorCtxKey).(*message.Printer)
	return mp, ok
}

// Wrapped responses

// OkRes builds an OK response including data and cero, one  or more messages.
// All messages are assumed of type info therefore flashes will be also of this type.
func (ep *Endpoint) OKRes(r *http.Request, data interface{}, msgs ...string) WrappedRes {
	fls := []FlashData{}
	for _, m := range msgs {
		fls = append(fls, ep.MakeFlash(m, InfoMT))
	}

	wr := WrappedRes{
		Data:  data,
		Flash: fls,
		Err:   nil,
	}

	wr.AddCSRF(r)

	return wr
}

// MultiRes builds a multiType response including data and cero, one  or more messages.
// Generated flash messages will be created according to the values passed in the parameter map.
// i.e.: map[string]MsgType{"Action processed": InfoMT, "Remember to update profile": WarnMT}
func (ep *Endpoint) MultiRes(r *http.Request, data interface{}, msgs map[string]MsgType) WrappedRes {
	fls := []FlashData{}
	for m, t := range msgs {
		fls = append(fls, ep.MakeFlash(m, t))
	}

	wr := WrappedRes{
		Data:  data,
		Flash: fls,
		Err:   nil,
	}

	wr.AddCSRF(r)

	return wr
}

// ErrRes builds an error response including data and cero, one  or more messages.
// All messages are assumed of type error therefore flashes will be also of this type.
func (ep *Endpoint) ErrRes(r *http.Request, data interface{}, msgs ...string) WrappedRes {
	fls := []FlashData{}
	for _, m := range msgs {
		fls = append(fls, ep.MakeFlash(m, ErrorMT))
	}

	wr := WrappedRes{
		Data:  data,
		Flash: fls,
		Err:   nil,
	}

	wr.AddCSRF(r)

	return wr
}

// Wrap response data.
func (ep *Endpoint) Wrap(r *http.Request, data interface{}, msg string, msgType MsgType, err error) WrappedRes {
	wr := WrappedRes{
		Data:  data,
		Flash: []FlashData{ep.MakeFlash(msg, msgType)},
		Err:   err,
	}

	wr.AddCSRF(r)

	return wr
}

// MakeFlash message.
func (ep *Endpoint) MakeFlash(msg string, msgType MsgType) FlashData {
	return FlashData{
		Type: msgType,
		Msg:  msg,
	}
}

// Redirect to url.
func (ep *Endpoint) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 302)
}

func (wr *WrappedRes) AddCSRF(r *http.Request) {
	wr.CSRF = map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}
}

// Forms
func FormToModel(r *http.Request, model interface{}) error {
	return NewDecoder().Decode(model, r.Form)
}

// NewDecoder build a schema decoder
// that put values from a map[string][]string into a struct.
func NewDecoder() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return d
}

// Resource paths

// Resource path functions
// IndexPath returns index path under resource root path.
func IndexPath() string {
	return ""
}

// EditPath returns edit path under resource root path.
func EditPath() string {
	return "/{id}/edit"
}

// NewPath returns new path under resource root path.
func NewPath() string {
	return "/new"
}

// ShowPath returns show path under resource root path.
func ShowPath() string {
	return "/{id}"
}

// CreatePath returns create path under resource root path.
func CreatePath() string {
	return ""
}

// UpdatePath returns update path under resource root path.
func UpdatePath() string {
	return "/{id}"
}

// InitDeletePath returns init delete path under resource root path.
func InitDeletePath() string {
	return "/{id}/init-delete"
}

// DeletePath returns delete path under resource root path.
func DeletePath() string {
	return "/{id}"
}

// SignupPath returns signup path.
func SignupPath() string {
	return "/signup"
}

// LoginPath returns login path.
func LoginPath() string {
	return "/login"
}

// ResPath
func ResPath(rootPath string) string {
	return "/" + rootPath + IndexPath()
}

// ResPathEdit
func ResPathEdit(rootPath string, r Identifiable) string {
	return fmt.Sprintf("/%s/%s/edit", rootPath, r.GetSlug())
}

// ResPathNew
func ResPathNew(rootPath string) string {
	return fmt.Sprintf("/%s/new", rootPath)
}

// ResPathInitDelete
func ResPathInitDelete(rootPath string, r Identifiable) string {
	return fmt.Sprintf("/%s/%s/init-delete", rootPath, r.GetSlug())
}

// ResPathSlug
func ResPathSlug(rootPath string, r Identifiable) string {
	return fmt.Sprintf("/%s/%s", rootPath, r.GetSlug())
}

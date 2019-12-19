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

	// "github.com/davecgh/go-spew/spew"

	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/markbates/pkger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/log"
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
	Localizer struct {
		*i18n.Localizer
	}
)

type (
	// WrappedRes stands for wrapped response
	WrappedRes struct {
		Data  interface{}
		Loc   *Localizer
		Flash FlashSet
		CSRF  map[string]interface{}
		Err   error
	}

	FlashSet []FlashItem

	// Flash data to present in page
	FlashItem struct {
		Msg  string
		Type MsgType
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
	NewTmpl     = "new.tmpl"
	IndexTmpl   = "index.tmpl"
	EditTmpl    = "edit.tmpl"
	ShowTmpl    = "show.tmpl"
	InitDelTmpl = "initdel.tmpl"
	SignUpTmpl  = "signup.tmpl"
	SignInTmpl  = "signin.tmpl"
)

const (
	InfoMT  MsgType = "info"
	WarnMT  MsgType = "warn"
	ErrorMT MsgType = "error"
	DebugMT MsgType = "debug"
)

var (
	InfoMTColor    = []string{"green-800", "white", "green-500", "green-800"}
	WarnMTColor    = []string{"yellow-800", "white", "yellow-500", "yellow-800"}
	ErrorMTColor   = []string{"red-800", "white", "red-500", "red-800"}
	DebugMTColor   = []string{"blue-800", "white", "blue-500", "blue-800"}
	DefaultMTColor = []string{"white", "white", "white", "white"}
)

const (
	FlashStoreKey = "flash"
)

const (
	I18NorCtxKey ContextKey = "i18n"
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

func (wr *WrappedRes) AllFlashes() FlashSet {
	fs := MakeFlashSet()

	for _, fi := range wr.Flash {
		fs = append(fs, fi)
	}

	return fs
}

func (l *Localizer) Localize(textID string) string {
	if l.Localizer != nil {
		t, _, err := l.LocalizeWithTag(&i18n.LocalizeConfig{
			MessageID: textID,
		})

		if err != nil {
			return fmt.Sprintf("%s", textID) // "'%s' [untransalted]", textID
		}

		return t
	}

	return fmt.Sprintf("%s", textID) // "'%s' [untransalted]", textID
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

			// ep.Log().Warn("Not a valid template", "path", path)

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

// Sessions
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

	ep.store = sessions.NewCookieStore([]byte(k))
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

func (ep *Endpoint) GetSession(r *http.Request, name ...string) *sessions.Session {
	session := "session"
	if len(name) > 0 {
		session = name[0]
	}
	s, err := ep.Store().Get(r, session)
	if err != nil {
		ep.Log().Warn("Cannot get sesssion from store", "reqID", "n/a")
	}
	return s
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

func getI18NLocalizer(r *http.Request) (localizer *i18n.Localizer, ok bool) {
	localizer, ok = r.Context().Value(I18NorCtxKey).(*i18n.Localizer)
	return localizer, ok
}

// Wrapped responses
// StdRes builds a multiType response including data, and info message and zero or more warning messages.
// If there are pending flash messages to show stored by previous action before redirecting
// they are also appended to the wrapped response.
func (ep *Endpoint) OKRes(w http.ResponseWriter, r *http.Request, data interface{}, infoMsg string, warnMsgs ...string) WrappedRes {
	f := MakeFlashSet()

	// Add info message
	m := strings.Trim(infoMsg, " ")
	if m != "" {
		f = f.AddItem(ep.MakeFlashItem(m, InfoMT))
	}

	// Add warnings
	if len(warnMsgs) > 0 {
		for _, m := range warnMsgs {
			f = f.AddItem(ep.MakeFlashItem(m, WarnMT))
		}
	}

	// Add pending messages
	f = f.AddItems(ep.RestoreFlash(r))

	wr := WrappedRes{
		Data:  data,
		Loc:   ep.Localizer(r),
		Flash: f,
		Err:   nil,
	}

	wr.AddCSRF(r)

	ep.ClearFlash(w, r)

	return wr
}

// ErrRes builds an error response including data and one error message.
// If there are pending flash messages to show stored by previous action before redirecting
// they are also appended to the wrapped response.
func (ep *Endpoint) ErrRes(w http.ResponseWriter, r *http.Request, data interface{}, errorMsg string, err error) WrappedRes {
	f := MakeFlashSet()

	// Add error message
	m := strings.Trim(errorMsg, " ")
	if m != "" {
		f = f.AddItem(ep.MakeFlashItem(m, ErrorMT))
	}

	// Add pending messages
	f = f.AddItems(ep.RestoreFlash(r))

	wr := WrappedRes{
		Data:  data,
		Loc:   ep.Localizer(r),
		Flash: f,
		Err:   err,
	}

	wr.AddCSRF(r)

	ep.ClearFlash(w, r)

	return wr
}

// Wrap response data.
func (ep *Endpoint) Wrap(r *http.Request, data interface{}, msg string, msgType MsgType, err error) WrappedRes {
	wr := WrappedRes{
		Data:  data,
		Loc:   ep.Localizer(r),
		Flash: FlashSet{ep.MakeFlashItem(msg, msgType)},
		Err:   err,
	}

	wr.AddCSRF(r)

	return wr
}

// Flash

// MakeFlash message.
func (ep *Endpoint) MakeFlashItem(msg string, msgType MsgType) FlashItem {
	return FlashItem{
		Msg:  msg,
		Type: msgType,
	}
}

func (ep *Endpoint) StoreFlash(w http.ResponseWriter, r *http.Request, message string, mt MsgType) (ok bool) {
	s := ep.GetSession(r)

	// Append to current ones
	f := ep.RestoreFlash(r)
	f = append(f, ep.MakeFlashItem(message, mt))

	s.Values[FlashStoreKey] = f
	err := s.Save(r, w)
	if err != nil {
		ep.Log().Error(err)
		return true
	}

	return false
}

func (ep *Endpoint) RestoreFlash(r *http.Request) FlashSet {
	s := ep.GetSession(r)
	v := s.Values[FlashStoreKey]

	f, ok := v.(FlashSet)
	if ok {
		//ep.Log().Debug("Stored flash", "value", spew.Sdump(f))
		return f
	}

	ep.Log().Info("No stored flash", "key", FlashStoreKey)
	return MakeFlashSet()
}

func (ep *Endpoint) ClearFlash(w http.ResponseWriter, r *http.Request) (ok bool) {
	s := ep.GetSession(r)
	delete(s.Values, FlashStoreKey)
	err := s.Save(r, w)
	if err != nil {
		return true
	}
	return false
}

// Redirect to url.
func (ep *Endpoint) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 302)
}

func (ep *Endpoint) RedirectWithFlash(w http.ResponseWriter, r *http.Request, url string, msg string, msgType MsgType) {
	ep.StoreFlash(w, r, msg, msgType)
	http.Redirect(w, r, url, 302)
}

func (ep *Endpoint) Localizer(r *http.Request) *Localizer {
	l, ok := getI18NLocalizer(r)
	if !ok {
		return nil
	}

	return &Localizer{l}
}

func (wr *WrappedRes) AddLocalizer(l *Localizer) {
	wr.Loc = l
}

func (wr *WrappedRes) AddCSRF(r *http.Request) {
	wr.CSRF = map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}
}

func (fi *FlashItem) Color() []string {
	if fi.Type == InfoMT {
		return InfoMTColor

	} else if fi.Type == WarnMT {
		return WarnMTColor

	} else if fi.Type == ErrorMT {
		return ErrorMTColor

	} else if fi.Type == DebugMT {
		return DebugMTColor

	}
	return DefaultMTColor
}

// Forms
func (ep *Endpoint) FormToModel(r *http.Request, model interface{}) error {
	return NewDecoder().Decode(model, r.Form)
}

// NewDecoder build a schema decoder
// that put values from a map[string][]string into a struct.
func NewDecoder() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return d
}

// Flash
func MakeFlashSet() FlashSet {
	return make(FlashSet, 0)
}

func (f FlashSet) IsEmpty() bool {
	return len(f) == 0
}

func (f FlashSet) AddItem(fi FlashItem) FlashSet {
	return append(f, fi)
}

func (f FlashSet) AddItems(fis []FlashItem) FlashSet {
	return append(f, fis...)
}

func (fi FlashItem) IsEmpty() bool {
	return fi == FlashItem{}
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

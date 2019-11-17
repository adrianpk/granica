package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/pkger"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
	"gitlab.com/mikrowezel/backend/log"
)

type (
	Endpoint struct {
		ctx         context.Context
		cfg         *config.Config
		log         *log.Logger
		service     *service.Service
		templates   TemplateSet
		templatesFx template.FuncMap
		store       *sessions.CookieStore
		storeKey    string
	}

	TemplateSet    map[string]*template.Template
	TemplateGroups map[string]map[string][]string
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

type (
	Identifiable interface {
		GetSlug() string
	}
)

const (
	templateDir = "/assets/web/embed/template"
	layoutDir   = "layout"
	layoutKey   = "layout"
	pageKey     = "page"
	partialKey  = "partial"
)

const (
	indexTmpl = "index.tmpl"
)

const (
	InfoMT  MsgType = "info"
	WarnMT  MsgType = "warn"
	ErrorMT MsgType = "error"
	DebugMT MsgType = "debug"
)

func MakeEndpoint(ctx context.Context, cfg *config.Config, log *log.Logger, service *service.Service) (*Endpoint, error) {
	e := Endpoint{
		ctx:     ctx,
		cfg:     cfg,
		log:     log,
		service: service,
	}

	// Cookie store
	e.makeCookieStore()

	// Load
	ts, err := e.loadTemplates()
	if err != nil {
		return &e, err
	}

	// Classify
	tg := e.classifyTemplates(ts)

	// Parse
	e.parseTemplates(ts, tg)

	return &e, nil
}

func (e *Endpoint) Ctx() context.Context {
	return e.ctx
}

func (e *Endpoint) Cfg() *config.Config {
	return e.cfg
}

func (e *Endpoint) Log() *log.Logger {
	return e.log
}

// loadTemplates from embedded filesystem (pkger)
// under '/assets/web/embed/template'
func (e *Endpoint) loadTemplates() (TemplateSet, error) {
	tmpls := make(TemplateSet)

	err := pkger.Walk(templateDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				e.Log().Error(err, "msg", "Cannot load template", "path", path)
				return err
			}

			if filepath.Ext(path) == ".tmpl" {
				list := filepath.SplitList(path)
				base := fmt.Sprintf("%s:%s", list[0], templateDir)
				p, _ := filepath.Rel(base, path)

				e.Log().Info("Reading template", "path", p)

				tmpls[p] = template.New(base)
				return nil
			}

			//e.Log().Warn("Not a valid template", "path", path)

			return nil
		})

	if err != nil {
		e.Log().Error(err, "msg", "Cannot load templates", "path")
		return tmpls, err
	}

	return tmpls, nil
}

// classifyTemplates grouping them,
// first by type (layout, partial and page)
// and then by resource.
func (e *Endpoint) classifyTemplates(ts TemplateSet) TemplateGroups {
	all := make(TemplateGroups)
	last := ""
	keys := e.tmplsKeys(ts)

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

func (e *Endpoint) tmplsKeys(ts TemplateSet) []string {
	keys := make([]string, 0, len(ts))
	for k, _ := range ts {
		keys = append(keys, k)
	}
	return keys
}

// parseTemplates parses template sets for each resource.
func (e *Endpoint) parseTemplates(ts TemplateSet, tg TemplateGroups) {
	e.templates = make(TemplateSet)
	layout := tg[layoutDir][layoutKey][0]

	for k, ts := range tg {
		pages := ts[pageKey]
		partials := ts[partialKey]

		for _, t := range pages {
			if k != layoutDir {
				e.parseTemplate(t, partials, layout, e.templatesFx)
			}
		}
	}
}

func (e *Endpoint) parseTemplate(page string, partials []string, layout string, funcs template.FuncMap) {
	parse := "base.tmpl"
	all := make([]string, 10)
	all = append(all, page)
	all = append(all, partials...)
	all = append(all, layout)
	trimSlice(&all)

	t, err := template.New(parse).Funcs(funcs).ParseFiles(all...)
	if err != nil {
		e.Log().Error(err, "Error parsing template set", "page", page)
	}

	base := fmt.Sprintf(".%s", templateDir)
	p, _ := filepath.Rel(base, page)

	e.Log().Info("Parsed template set", "path", p)

	e.templates[page] = t
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
func (e *Endpoint) makeCookieStore() {
	k := e.Cfg().ValOrDef("web.cookiestore.key", "")
	if k == "" {
		k = e.genAES256Key()
		e.Log().Debug("New cookie store random key", "value", k)
		csEnVar := fmt.Sprintf("%s_COOKIESTORE_KEY", "GRN")
		e.Log().Info("Set a custom cookie store key using a 32 char string stored as an envar", "envvar", csEnVar)
	}

	e.storeKey = k
	e.Log().Debug("Cookie store key", "value", k)
	sessions.NewCookieStore([]byte(k))
}

func (e *Endpoint) genAES256Key() string {
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

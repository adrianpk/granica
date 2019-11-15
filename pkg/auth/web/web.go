package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/pkger"
	"gitlab.com/mikrowezel/backend/config"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/service"
	"gitlab.com/mikrowezel/backend/log"
)

type (
	Endpoint struct {
		ctx     context.Context
		cfg     *config.Config
		log     *log.Logger
		service *service.Service
		parsed  TemplateSet
	}

	TemplateSet    map[string]*template.Template
	TemplateGroups map[string]map[string][]string
)

type (
	contextKey string
)

type (
	Identifiable interface {
		GetSlug() string
	}
)

const (
	templateDir = "/assets/web/template"
	layoutDir   = "layout"
	layoutKey   = "layout"
	pageKey     = "page"
	partialKey  = "partial"
)

func MakeEndpoint(ctx context.Context, cfg *config.Config, log *log.Logger, service *service.Service) (*Endpoint, error) {
	e := Endpoint{
		ctx:     ctx,
		cfg:     cfg,
		log:     log,
		service: service,
	}

	ts, err := e.collectTemplates()
	if err != nil {
		return &e, err
	}
	tg := e.classifyTemplates(ts)
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

// collectTemplates embedded filesystem (pkger)
// under '/assets/web/template'
func (e *Endpoint) collectTemplates() (TemplateSet, error) {
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

				e.Log().Info("Template file", "path", p)

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

// classifyTemplates organizing them,
// first by type (layout, partial and page)
// and then by resource.
func (e *Endpoint) classifyTemplates(ts TemplateSet) TemplateGroups {
	all := make(TemplateGroups)
	last := ""
	keys := e.tmplsKeys(ts)

	for _, path := range keys {
		p := "./assets/web/template" + "/" + path

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

func (e *Endpoint) parseTemplates(ts TemplateSet, tg TemplateGroups) {
	e.parsed = make(TemplateSet)
	layout := tg[layoutDir][layoutKey][0]

	for k, ts := range tg {
		pages := ts[pageKey]
		partials := ts[partialKey]

		for _, t := range pages {
			if k != layoutDir {
				e.parseTemplate(t, partials, layout, pathFxs)
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

	e.Log().Info("Template processed", "template", page)

	e.parsed[page] = t
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

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
		ctx       context.Context
		cfg       *config.Config
		log       *log.Logger
		service   *service.Service
		templates TemplateSet
	}

	TemplateResSet map[string][]*template.Template
	TemplateSet    map[string]*template.Template
)

type (
	contextKey string
)

type (
	Identifiable interface {
		GetSlug() string
	}
)

func MakeEndpoint(ctx context.Context, cfg *config.Config, log *log.Logger, service *service.Service) *Endpoint {
	e := Endpoint{
		ctx:     ctx,
		cfg:     cfg,
		log:     log,
		service: service,
	}

	e.collectTemplates()

	return &e
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

func (e *Endpoint) collectTemplates() error {
	ts := TemplateSet{}

	err := pkger.Walk("/",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				e.Log().Error(err, "msg", "Cannot load template", "path", path)
				return err
			}

			if filepath.Ext(path) == ".tmpl" {
				e.Log().Info("New template file", "path", path)
				ts[path] = template.New(path)
			}
			return nil
		})

	if err != nil {
		e.Log().Error(err, "msg", "Cannot load templates", "path")
		return err
	}
	return nil
}

//func (e *Endpoint) loadTemplates() {
//ts := ac.Templates
//assets := bindata.Resource(tmpl.AssetNames(),
//func(name string) ([]byte, error) {
//return tmpl.Asset(name)
//})

//ct := boot.ClassifyTemplates(assets)
//layout := ct["layouts"]["app"][0]
//for k, ts := range ct {
//standard := ts["standard"]
//partials := ts["partials"]
//for _, tt := range standard {
//if k != "layouts" {
//ParseAsset(tt, partials, layout, common.Routes, templatesMap)
//}
//}
//}
//}

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

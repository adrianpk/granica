package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/markbates/pkger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/mikrowezel/backend/web"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type textResponse string

var (
	langMatcher = language.NewMatcher(message.DefaultCatalog.Languages())
)

func (t textResponse) write(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(t))
}

// AddWebServer to worker.
func (a *Auth) AddWebServer() http.Handler {
	// Home
	hr := a.makeHomeWebRouter()

	// User
	a.makeUserWebRouter(hr)

	// Account
	//a.makeAccountWebRouter(hr)

	a.WebServer = hr

	return hr
}

// AddJSONRESTServer to worker.
func (a *Auth) AddJSONRESTServer() http.Handler {
	// Home
	hr := a.makeHomeJSONRESTRouter()

	// API
	ar := a.makeAPIJSONRESTRouter(hr)

	// User
	a.makeUserJSONRESTRouter(ar)

	// Account
	a.makeAccountJSONRESTRouter(ar)

	a.JSONRESTServer = hr

	return hr
}

func (a *Auth) makeHomeWebRouter() chi.Router {
	hr := chi.NewRouter()
	hr.Use(middleware.RequestID)
	hr.Use(middleware.RealIP)
	hr.Use(middleware.Recoverer)
	hr.Use(middleware.Timeout(60 * time.Second))
	hr.Use(a.MethodOverride)
	hr.Use(a.CSRFProtection)
	hr.Use(a.I18N)
	a.addHomeWebRoutes(hr)
	return hr
}

func (a *Auth) makeHomeJSONRESTRouter() chi.Router {
	hr := chi.NewRouter()
	hr.Use(middleware.RequestID)
	hr.Use(middleware.RealIP)
	hr.Use(middleware.Recoverer)
	hr.Use(middleware.Timeout(60 * time.Second))
	a.addHomeJSONRESTRoutes(hr)
	return hr
}

func (a *Auth) addHomeWebRoutes(rt chi.Router) {
	dir := "/assets/web/embed/public"
	fs := http.FileServer(FileSystem{pkger.Dir(dir)})

	rt.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := pkger.Stat(dir + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)

		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

func (a *Auth) addHomeJSONRESTRoutes(rt chi.Router) {
	rt.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tr := textResponse("Granica JSON REST API is running!")
		rt.Get("/", tr.write)
	})
}

func (a *Auth) makeAPIJSONRESTRouter(parent chi.Router) chi.Router {
	return parent.Route("/api/v1", func(ar chi.Router) {
		tr := textResponse("API v1.0")
		ar.Get("/", tr.write)
	})
}

// Middleware

// MethodOverride to emulate PUT and PATCH HTTP methods.
func (a *Auth) MethodOverride(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			method := r.PostFormValue("_method")
			if method == "" {
				method = r.Header.Get("X-HTTP-Method-Override")
			}
			if method == "PUT" || method == "PATCH" || method == "DELETE" {
				r.Method = method
			}
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// CSRFProtection add cross-site request forgery protecction to the handler.
func (a *Auth) CSRFProtection(h http.Handler) http.Handler {
	return csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))(h)
}

// I18N

func (a *Auth) I18N(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		lang := r.FormValue("lang")
		accept := r.Header.Get("Accept-Language")
		bundle := a.I18NBundle()

		l := i18n.NewLocalizer(bundle, lang, accept)
		ctx := context.WithValue(r.Context(), web.I18NorCtxKey, l)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (a *Auth) I18NBundle() *i18n.Bundle {
	if a.i18nBundle != nil {
		return a.i18nBundle
	}

	b := i18n.NewBundle(language.English)
	b.RegisterUnmarshalFunc("json", json.Unmarshal)

	locales := []string{"en", "pl", "de", "es"}

	for _, loc := range locales {
		path := fmt.Sprintf("/assets/web/embed/i18n/%s.json", loc)

		// Open pkger filer
		f, err := pkger.Open(path)
		if err != nil {
			a.Log().Error(err, "path", path)
		}
		defer f.Close()

		// Read file content
		fs, err := f.Stat()
		if err != nil {
			a.Log().Error(err, "path", path)
		}
		fd := make([]byte, fs.Size())
		f.Read(fd)

		// Load into bundle
		b.ParseMessageFileBytes(fd, path)
		//b.LoadEmbeddedMessageFile(fd, path)
	}

	// Cache it
	a.i18nBundle = b

	return b
}

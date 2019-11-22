package auth

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/markbates/pkger"
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
	hr.Use(CSRFProtection)
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

// CSRFProtection add cross-site request forgery protecction to the handler.
func CSRFProtection(h http.Handler) http.Handler {
	return csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))(h)
}

// I18N
func I18N(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// NOTE: Lang is read from a query string.
		// TODO: Read lang from headers and/or value
		// stored in cookie or user session.
		l, ok := r.URL.Query()["lang"]

		if !ok || len(l) < 1 {
			l = append(l, language.English.String())
		}

		tag, _, _ := langMatcher.Match(language.MustParse(l[0]))

		p := message.NewPrinter(tag)
		ctx := context.WithValue(context.Background(), web.I18NorCtxKey, p)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

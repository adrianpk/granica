package auth

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type textResponse string

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
	rt.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tr := textResponse("Granica web server is running!")
		rt.Get("/", tr.write)
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

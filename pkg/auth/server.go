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

// AddServer to worker.
func (a *Auth) AddServer() http.Handler {
	// Home
	hr := a.makeHomeRouter()

	// API
	ar := a.makeAPIRouter(hr)

	// User
	a.makeUserAPIRouter(ar)

	// Account
	a.makeAccountAPIRouter(ar)

	a.Server = hr

	return hr
}

func (a *Auth) makeHomeRouter() chi.Router {
	hr := chi.NewRouter()
	hr.Use(middleware.RequestID)
	hr.Use(middleware.RealIP)
	hr.Use(middleware.Recoverer)
	hr.Use(middleware.Timeout(60 * time.Second))
	a.addHomeRoutes(hr)
	return hr
}

func (a *Auth) addHomeRoutes(rt chi.Router) {
	rt.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tr := textResponse("Granica is running!")
		rt.Get("/", tr.write)
	})
}

func (a *Auth) makeAPIRouter(parent chi.Router) chi.Router {
	return parent.Route("/api/v1", func(ar chi.Router) {
		tr := textResponse("API v1.0")
		ar.Get("/", tr.write)
	})
}

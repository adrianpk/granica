package auth

import (
	"context"
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
	// Middlewares
	hr := chi.NewRouter()
	hr.Use(middleware.RequestID)
	hr.Use(middleware.RealIP)
	hr.Use(middleware.Recoverer)
	hr.Use(middleware.Timeout(60 * time.Second))
	a.addHomeRoutes(hr)

	// API
	ar := a.makeAPIRouter(hr)

	// User
	a.makeUserAPIRouter(ar)

	a.Server = hr

	return hr
}

func (a *Auth) addHomeRoutes(rt chi.Router) {
	rt.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tr := textResponse("Granica is running!")
		rt.Get("/", tr.write)
	})
}

func (a *Auth) makeAPIRouter(parent chi.Router) chi.Router {
	return parent.Route("/api/v1", func(apiRt chi.Router) {
		tr := textResponse("API v1.0")
		apiRt.Get("/", tr.write)
	})
}

func (a *Auth) makeUserAPIRouter(parent chi.Router) chi.Router {
	return parent.Route("/users", func(userAPIRt chi.Router) {
		userAPIRt.Use(userCtx)
		userAPIRt.Get("/", a.getUsers)
	})
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		someVal := ""
		ctx := context.WithValue(r.Context(), userCtxKey, someVal)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

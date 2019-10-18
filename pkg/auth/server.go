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

// We usually use slug to avoid exposing database ID to the world.
// But because but as this value is immutable and includes the
// username selected when the user was created/registered we
// prefer to use username as the external main identifier.
func (a *Auth) makeUserAPIRouter(parent chi.Router) chi.Router {
	return parent.Route("/users", func(uar chi.Router) {
		uar.Post("/", a.CreateUserJSON)
		uar.Get("/", a.GetUsersJSON)
		uar.Route("/{username}", func(uarid chi.Router) {
			uarid.Use(userCtx)
			uarid.Get("/", a.GetUserJSON)
			uarid.Put("/", a.UpdateUserJSON)
			uarid.Delete("/", a.DeleteUserJSON)
		})
	})
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		ctx := context.WithValue(r.Context(), userCtxKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Account
func (a *Auth) makeAccountAPIRouter(parent chi.Router) chi.Router {
	return parent.Route("/accounts", func(aar chi.Router) {
		aar.Post("/", a.CreateAccountJSON)
		aar.Get("/", a.GetAccountJSON)
		aar.Route("/{account}", func(aarid chi.Router) {
			aarid.Use(accountCtx)
			aarid.Get("/", a.GetAccountJSON)
			aarid.Put("/", a.UpdateAccountJSON)
			aarid.Delete("/", a.DeleteAccountJSON)
		})
	})
}

func accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "account-slug")
		ctx := context.WithValue(r.Context(), accountCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

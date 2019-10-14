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
		uar.Post("/", a.createUserJSON)
		uar.Get("/", a.getUsersJSON)
		uar.Route("/{username}", func(uarid chi.Router) {
			uarid.Use(userCtx)
			uarid.Get("/", a.getUserJSON)
			uarid.Put("/", a.updateUserJSON)
			uarid.Delete("/", a.deleteUserJSON)
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

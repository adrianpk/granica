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
	rt := chi.NewRouter()
	rt.Use(middleware.RequestID)
	rt.Use(middleware.RealIP)
	rt.Use(middleware.Recoverer)
	rt.Use(middleware.Timeout(60 * time.Second))

	rt.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tr := textResponse("Granica is running!")
		rt.Get("/", tr.write)
	})

	apiRt := rt.Route("/api/v1", func(apiRt chi.Router) {
		tr := textResponse("API v1.0")
		apiRt.Get("/", tr.write)
	})

	apiRt.Route("/users", func(userRt chi.Router) {
		userRt.Use(userCtx)
		userRt.Get("/", a.getUsers)
	})

	a.Server = rt

	return rt
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		someVal := ""
		ctx := context.WithValue(r.Context(), userCtxKey, someVal)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

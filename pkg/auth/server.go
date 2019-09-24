package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// MakeServer creates worker server.
func MakeServer() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Granica is running!"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(m.authCtx)
		r.Get("/", m.getUsersHandler) // POST /routes
	})

	return r
}

func (m *Manager) userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		someVal := ""
		ctx := context.WithValue(r.Context(), userCtxKey, someVal)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

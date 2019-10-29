package auth

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/jsonrest"
)

// Account
func (a *Auth) makeAccountAPIRouter(parent chi.Router) chi.Router {
	return parent.Route("/accounts", func(aar chi.Router) {
		aar.Post("/", a.jsonep.CreateAccount)
		aar.Get("/", a.jsonep.GetAccount)
		aar.Route("/{account}", func(aarid chi.Router) {
			aarid.Use(accountCtx)
			aarid.Get("/", a.jsonep.GetAccount)
			aarid.Put("/", a.jsonep.UpdateAccount)
			aarid.Delete("/", a.jsonep.DeleteAccount)
		})
	})
}

func accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "account-slug")
		ctx := context.WithValue(r.Context(), jsonrest.AccountCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

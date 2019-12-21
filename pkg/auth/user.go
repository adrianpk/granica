package auth

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/web"
)

func (a *Auth) makeUserWebRouter(parent chi.Router) chi.Router {
	return parent.Route("/users", func(uar chi.Router) {
		uar.Get("/", a.webep.IndexUsers)
		uar.Get("/new", a.webep.NewUser)
		uar.Post("/", a.webep.CreateUser)
		uar.Get("/signup", a.webep.InitSignUpUser)
		uar.Post("/signup", a.webep.SignUpUser)
		uar.Get("/signin", a.webep.InitSignInUser)
		uar.Post("/signin", a.webep.SignInUser)
		uar.Route("/{slug}", func(uarid chi.Router) {
			uarid.Use(userCtx)
			uarid.Get("/", a.webep.ShowUser)
			uarid.Get("/edit", a.webep.EditUser)
			uarid.Patch("/", a.webep.UpdateUser)
			uarid.Put("/", a.webep.UpdateUser)
			uarid.Post("/init-delete", a.webep.InitDeleteUser)
			uarid.Delete("/", a.webep.DeleteUser)
			uarid.Route("/{token}", func(uartkn chi.Router) {
				uartkn.Use(confCtx)
				uartkn.Get("/confirm", a.webep.ConfirmUser)
			})
		})
	})
}

func (a *Auth) makeUserJSONRESTRouter(parent chi.Router) chi.Router {
	return parent.Route("/users", func(uar chi.Router) {
		uar.Post("/", a.jsonep.CreateUser)
		uar.Get("/", a.jsonep.IndexUsers)
		uar.Route("/{slug}", func(uarid chi.Router) {
			uarid.Use(userCtx)
			uarid.Get("/", a.jsonep.GetUser)
			uarid.Patch("/", a.jsonep.UpdateUser)
			uarid.Put("/", a.jsonep.UpdateUser)
			uarid.Delete("/", a.jsonep.DeleteUser)
		})
	})
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		ctx := context.WithValue(r.Context(), web.UserCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func confCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "token")
		ctx := context.WithValue(r.Context(), web.ConfCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

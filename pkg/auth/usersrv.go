package auth

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/mikrowezel/backend/granica/pkg/auth/jsonrest"
)

func (a *Auth) makeUserWebRouter(parent chi.Router) chi.Router {
	return parent.Route("/users", func(uar chi.Router) {
		//uar.Post("/", a.webep.CreateUser)
		uar.Get("/", a.webep.GetUsers)
		uar.Route("/{username}", func(uarid chi.Router) {
			uarid.Use(userCtx)
			//uarid.Get("/", a.jsonep.GetUser)
			//uarid.Put("/", a.jsonep.UpdateUser)
			//uarid.Delete("/", a.jsonep.DeleteUser)
		})
	})
}

// We usually use slug to avoid exposing database ID to the world.
// But because but as this value is immutable and includes the
// username selected when the user was created/registered we
// prefer to use username as the external main identifier.
func (a *Auth) makeUserJSONRESTRouter(parent chi.Router) chi.Router {
	return parent.Route("/users", func(uar chi.Router) {
		uar.Post("/", a.jsonep.CreateUser)
		uar.Get("/", a.jsonep.GetUsers)
		uar.Route("/{username}", func(uarid chi.Router) {
			uarid.Use(userCtx)
			uarid.Get("/", a.jsonep.GetUser)
			uarid.Put("/", a.jsonep.UpdateUser)
			uarid.Delete("/", a.jsonep.DeleteUser)
		})
	})
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		ctx := context.WithValue(r.Context(), jsonrest.UserCtxKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

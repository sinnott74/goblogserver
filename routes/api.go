package routes

import (
	"github.com/go-chi/chi"
	"github.com/sinnott74/goblogserver/middleware"
)

// APIRouter create all API routes
func APIRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Transaction)
	r.Mount("/users", UserRouter())
	r.Mount("/auth", AuthRouter())
	r.Mount("/blogposts", BlogpostRouter())
	return r
}

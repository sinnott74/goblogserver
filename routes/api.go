package routes

import (
	"github.com/go-chi/chi"
)

func ApiRouter() chi.Router {
	r := chi.NewRouter()
	r.Mount("/helloworld", HelloWorldRouter())
	r.Mount("/users", UserRouter())
	r.Mount("/auth", AuthRouter())
	r.Mount("/blogposts", BlogpostRouter())
	return r
}

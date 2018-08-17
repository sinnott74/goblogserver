package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// HelloWorldRouter creates the hello world routes
func HelloWorldRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getHelloWorld)
	return r
}

func getHelloWorld(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "Hello world")
}

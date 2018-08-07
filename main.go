package main

import (
	"net/http"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/sinnott74/goblogserver/database"
	"github.com/sinnott74/goblogserver/env"
	"github.com/sinnott74/goblogserver/middleware"
	"github.com/sinnott74/goblogserver/routes"
)

func main() {

	err := database.Init()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.ForceHTTPS)
	r.Use(chiMiddleware.StripSlashes)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.DefaultCompress)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.Transaction)
	r.Mount("/api", routes.ApiRouter())

	http.ListenAndServe(":"+env.Port(), r)
}

package routes

import (
	"net/http"
	"net/http/pprof"

	"github.com/sinnott74/goblogserver/orm"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/sinnott74/goblogserver/env"
	"github.com/sinnott74/goblogserver/middleware"
)

// Handler creates all the middleware & routes for the http server
func Handler() http.Handler {

	debug := env.Debug()

	r := chi.NewRouter()
	r.Use(middleware.ForceHTTPS)
	if debug {
		r.Use(chiMiddleware.Logger)
	}
	r.Use(chiMiddleware.StripSlashes)
	r.Use(chiMiddleware.DefaultCompress)
	r.Use(middleware.Etag)
	r.Use(middleware.Recoverer)
	r.Mount("/api", APIRouter())
	r.Mount("/helloworld", HelloWorldRouter())

	// Debugging
	if debug {
		// Register pprof handlers
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	MapErrorToResponseCodes()

	return r
}

func MapErrorToResponseCodes() {
	middleware.MapErrorResponseStatus(&orm.NotFoundError{}, 404)
}
